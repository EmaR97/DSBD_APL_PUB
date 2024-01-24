from flask import Flask, request, jsonify, abort

from mongo import SLADocument
from prometheus.query import PrometheusInterface


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

    current_value = SlaManger.prometheus.get_prometheus_instant(sla.metric_name)
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
    values = SlaManger.prometheus.get_prometheus_vector(metric_name, start_time)
    sla = SLADocument.objects(metric_name=metric_name).first()
    violations = [value for value in values if sla.get_sla_status(value)]
    return jsonify({"violations_count": len(violations)})


# API Endpoint to Query Probability of Violations
@SlaManger.app.route('/probability', methods=['GET'])
def query_probability():
    pass  # metric_name = request.args.get('metric_name')  # x_minutes = int(request.args.get('minutes', 5))  #  # #
    # Simulate calculating the probability (replace with actual logic)  # probability = random.uniform(0,
    # 1)  #  # return jsonify({"probability": probability})
