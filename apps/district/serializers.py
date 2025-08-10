from rest_framework import serializers
from .models import District
from apps.city.serializers import CitySerializer


class DistrictSerializer(serializers.ModelSerializer):
    city = CitySerializer(read_only=True)
    """
    Serializador para o modelo District (Bairro).

    Campos:
    - id: Identificador único do bairro
    - name: Nome do bairro
    - city: Chave estrangeira para a cidade à qual o bairro pertence
    """
    class Meta:
        model = District
        fields = ['id', 'name', 'city']
        swagger_schema_fields = {
            "title": "Bairro",
            "description": "Representa um bairro de uma cidade"
        }
