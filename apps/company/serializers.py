from rest_framework import serializers
from .models import Company


class CompanySerializer(serializers.ModelSerializer):
    """
    Serializador para o modelo Company (Empresa).

    Campos:
    - id: Identificador único da empresa
    - cnpj: CNPJ da empresa (8 dígitos)
    - social_name: Razão social da empresa
    - juridical_nature: Natureza jurídica
    - responsible_qualification: Qualificação do responsável
    - social_capital: Capital social
    - company_size: Porte da empresa
    - federative_entity: Ente federativo
    - state: Estado onde a empresa está localizada
    - city: Cidade onde a empresa está localizada
    - district: Bairro onde a empresa está localizada
    - created_at: Data de criação do registro
    - updated_at: Data da última atualização
    """
    class Meta:
        model = Company
        fields = [
            'id', 'cnpj', 'social_name', 'juridical_nature',
            'responsible_qualification', 'social_capital', 'company_size',
            'federative_entity', 'state', 'city', 'district',
            'created_at', 'updated_at'
        ]
        swagger_schema_fields = {
            "title": "Empresa",
            "description": "Representa uma empresa e suas informações cadastrais"
        }
