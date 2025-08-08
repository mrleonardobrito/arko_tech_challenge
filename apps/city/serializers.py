from rest_framework import serializers
from .models import City


class CitySerializer(serializers.ModelSerializer):
    """
    Serializador para o modelo City (Cidade).

    Campos:
    - id: Identificador único da cidade
    - name: Nome da cidade
    - state: Chave estrangeira para o estado ao qual a cidade pertence
    - created_at: Data de criação do registro
    - updated_at: Data da última atualização
    """
    class Meta:
        model = City
        fields = ['id', 'name', 'state', 'created_at', 'updated_at']
        swagger_schema_fields = {
            "title": "Cidade",
            "description": "Representa uma cidade brasileira e seu estado"
        }
