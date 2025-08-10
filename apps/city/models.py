from django.db import models


class City(models.Model):
    name = models.CharField(max_length=255)
    state = models.ForeignKey('state.State', on_delete=models.CASCADE)

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = 'Cidade'
        verbose_name_plural = 'Cidades'
        ordering = ['name']
        db_table = 'city'
