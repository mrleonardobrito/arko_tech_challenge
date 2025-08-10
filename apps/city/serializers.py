from rest_framework import serializers
from .models import City
from apps.state.serializers import StateSerializer


class CitySerializer(serializers.ModelSerializer):
    state = StateSerializer(read_only=True)
    """
    Serializador para o modelo City (Cidade).

    Campos:
    - id: Identificador único da cidade
    - name: Nome da cidade
    - state: Chave estrangeira para o estado ao qual a cidade pertence
    """
    class Meta:
        model = City
        fields = ['id', 'name', 'state']
        swagger_schema_fields = {
            "title": "Cidade",
            "description": "Representa uma cidade brasileira e seu estado"
        }
