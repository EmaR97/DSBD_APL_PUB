import datetime

from mongoengine import Document, StringField, FloatField, DateTimeField


# MongoDB Document Model for SLA
class SLADocument(Document):
    metric_name = StringField(required=True)
    range_min = FloatField(required=True)
    range_max = FloatField(required=True)
    created_at = DateTimeField(default=datetime.datetime.now(datetime.UTC))

    # Helper function for query_sla
    def get_sla_status(self, current_value):
        return self.range_min > float(current_value) or self.range_max < float(current_value)
