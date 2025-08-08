from django.db import models


class City(models.Model):
    name = models.CharField(max_length=255)
    state = models.ForeignKey('state.State', on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = 'Cidade'
        verbose_name_plural = 'Cidades'
        ordering = ['name']
        db_table = 'city'
        unique_together = ['name', 'state']
        constraints = [
            models.UniqueConstraint(
                fields=['name', 'state'], name='unique_city_state')
        ]
