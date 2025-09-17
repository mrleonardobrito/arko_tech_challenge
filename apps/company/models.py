from django.db import models
# Create your models here.


class Company(models.Model):
    cnpj = models.CharField(max_length=8, unique=True, null=True)
    social_name = models.CharField(max_length=255, null=True)
    juridical_nature = models.CharField(max_length=4, null=True)
    responsible_qualification = models.CharField(max_length=2, null=True)
    social_capital = models.DecimalField(
        max_digits=15, decimal_places=2, null=True)
    company_size = models.CharField(max_length=2, null=True)
    federative_entity = models.CharField(max_length=100, null=True)

    def __str__(self):
        return self.social_name

    class Meta:
        verbose_name = 'Empresa'
        verbose_name_plural = 'Empresas'
        ordering = ['social_name']
        db_table = 'company'
        unique_together = ['cnpj']
        constraints = [
            models.UniqueConstraint(fields=['cnpj'], name='unique_cnpj')
        ]
        indexes = [
            models.Index(fields=['social_name'],
                         name='company_social_name_idx'),
            models.Index(fields=['cnpj'], name='company_cnpj_idx'),
            models.Index(fields=['federative_entity'],
                         name='company_federative_entity_idx'),
            models.Index(fields=['social_capital'],
                         name='company_social_capital_idx'),
            models.Index(fields=['company_size'],
                         name='company_company_size_idx'),
        ]
