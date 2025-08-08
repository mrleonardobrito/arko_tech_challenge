from django.db import models

# Create your models here.


class Company(models.Model):
    cnpj = models.CharField(max_length=8, unique=True)
    social_name = models.CharField(max_length=255)
    juridical_nature = models.CharField(max_length=4)
    responsible_qualification = models.CharField(max_length=2)
    social_capital = models.DecimalField(max_digits=15, decimal_places=2)
    company_size = models.CharField(max_length=2)
    federative_entity = models.CharField(max_length=4)
    state = models.ForeignKey('state.State', on_delete=models.CASCADE)
    city = models.ForeignKey('city.City', on_delete=models.CASCADE)
    district = models.ForeignKey('district.District', on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

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
