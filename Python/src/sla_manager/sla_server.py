from datetime import datetime, timedelta

from flask import Flask, request, jsonify, abort

from mongo.series_model import SeriesModel
from mongo.sla_document import SLADocument
from prometheus.query import PrometheusInterface
from time_series import compute_out_of_bounds_probability
from time_series import reevaluate_model


class SlaManger:
    app = Flask(__name__)
    prometheus: PrometheusInterface = None


# API Endpoint to Create/Update SLA
@SlaManger.app.route('/sla', methods=['PUT'])
def create_update_sla():
    data = request.json

    metric_name = data.get('metric_name')
    range_min = data.get('range_min')
    range_max = data.get('range_max')

    sla = SLADocument.objects(metric_name=metric_name).first()

    if sla:
        try:
            sla.update(set__range_min=range_min, set__range_max=range_max)
        except Exception as e:
            SlaManger.app.logger.error(f"Error updating SLA: {e}")
            abort(500, description="Internal Server Error")
    else:
        try:
            SLADocument(metric_name=metric_name, range_min=range_min, range_max=range_max).save()
        except Exception as e:
            SlaManger.app.logger.error(f"Error creating SLA: {e}")
            abort(500, description="Internal Server Error")

    return jsonify({"message": "SLA created/updated successfully"}), 201


# API Endpoint to Query SLA
@SlaManger.app.route('/sla', methods=['GET'])
def query_sla():
    metric_name = request.args.get('metric_name')
    sla = SLADocument.objects(metric_name=metric_name).first()

    if not sla:
        abort(404, description="SLA not found")

    current_value = SlaManger.prometheus.get_prometheus_instant(sla.metric_name)[1]
    violation = sla.get_sla_status(current_value)
    return jsonify({"metric_name": sla.metric_name, "current_value": current_value, "violation": violation,
                    "created_at": sla.created_at})


# API Endpoint to Remove SLA
@SlaManger.app.route('/sla', methods=['DELETE'])
def remove_sla():
    metric_name = request.args.get('metric_name')
    sla = SLADocument.objects(metric_name=metric_name).first()

    if not sla:
        abort(404, description="SLA not found")

    try:
        sla.delete()
    except Exception as e:
        SlaManger.app.logger.error(f"Error deleting SLA: {e}")
        abort(500, description="Internal Server Error")

    return jsonify({"message": "SLA deleted successfully"}), 200


# API Endpoint to Query Violations
@SlaManger.app.route('/violations', methods=['GET'])
def query_violations():
    metric_name = request.args.get('metric_name')
    start_time = request.args.get('hours')
    values = [x[1] for x in SlaManger.prometheus.get_prometheus_vector(metric_name, start_time)]
    sla = SLADocument.objects(metric_name=metric_name).first()
    if not sla:
        abort(404, description="SLA not found")
    violations = [value for value in values if sla.get_sla_status(value)]
    return jsonify({"violations_count": len(violations)})


# API Endpoint to Query Probability of Violations
@SlaManger.app.route('/probability', methods=['GET'])
def query_probability():
    # Retrieve parameters from the request
    metric_name = request.args.get('metric_name')
    x_minutes = int(request.args.get('minutes', 5))

    # Fetch SLA Document and necessary parameters
    sla = SLADocument.objects(metric_name=metric_name).first()
    if not sla:
        return jsonify({"error": f"No SLA Document found with metric name '{metric_name}'"}), 404

    y_lower_bound = sla.range_max

    # Define time boundaries
    current_datetime = datetime.now()
    x_lower_limit = int(current_datetime.timestamp())
    future_datetime = current_datetime + timedelta(minutes=x_minutes)
    x_upper_limit = int(future_datetime.timestamp())

    # Fetch SeriesModel instance and necessary parameters
    model_instance = SeriesModel.objects(metric_name=metric_name).first()
    if not model_instance:
        return jsonify({"error": f"No SeriesModel found with metric name '{metric_name}'"}), 404

    trend_function = model_instance.get_trend()
    error_std = model_instance.error_std

    try:
        # Compute the probability of y > y_lower_bound
        probability, _, _, _ = compute_out_of_bounds_probability(trend_function, error_std, x_lower_limit,
                                                                 x_upper_limit, y_lower_bound)
        # Log the result
        SlaManger.app.logger.info(
            f"Probability of y > {y_lower_bound} for {x_lower_limit} < x < {x_upper_limit} is {probability}%")

        # Return the result to the client
        return jsonify({"probability": probability})
    except Exception as e:
        # Log the error
        SlaManger.app.logger.error(f"An error occurred: {str(e)}")
        # Return an error response to the client
        return jsonify({"error": "An unexpected error occurred"}), 500


# API Endpoint to re-evaluate model for a specified metric
@SlaManger.app.route('/reevaluate_model', methods=['POST'])
def reevaluate_model_endpoint():
    data = request.json
    metric_name = data.get('metric_name')
    range_in_minute = data.get('range_in_minute')

    if not metric_name or not range_in_minute:
        abort(400)

    # Query Prometheus for time series data
    # Adjust this according to your PrometheusInterface implementation
    time_series_data = SlaManger.prometheus.get_prometheus_vector(metric_name, range_in_minute)

    if time_series_data is None:
        abort(404, description=f"No time series data found for metric '{metric_name}'")

    # Re-evaluate the model and get trend function and error_std
    trend_function, error_std = reevaluate_model(time_series_data)

    # Store the new model in the database
    try:
        model_instance = SeriesModel.objects(metric_name=metric_name).first()
        if model_instance:
            model_instance.update(set__error_std=error_std)
            model_instance.set_trend(trend_function)
        else:
            SeriesModel(metric_name=metric_name, error_std=error_std).save()
            model_instance = SeriesModel.objects(metric_name=metric_name).first()
            model_instance.set_trend(trend_function)
        model_instance.save()
    except Exception as e:
        abort(500, description=f"Error storing the new model for metric '{metric_name}': {str(e)}")

    return jsonify({"message": f"Model re-evaluation completed for metric '{metric_name}'"}), 200
