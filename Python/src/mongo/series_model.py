import dill
import mongoengine as me


# Define a model to store std, serialized trend data, and metric name
class SeriesModel(me.Document):
    error_std = me.FloatField(required=True)
    serialized_trend = me.BinaryField(required=True)
    metric_name = me.StringField(required=True, unique=True)

    def set_trend(self, trend_func):
        self.serialized_trend = dill.dumps(trend_func)

    def get_trend(self):
        if self.serialized_trend:
            return dill.loads(self.serialized_trend)
        else:
            return None

