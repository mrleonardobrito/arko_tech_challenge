from django.db import models

# Create your models here.


class District(models.Model):
    name = models.CharField(max_length=255)
    city = models.ForeignKey('city.City', on_delete=models.CASCADE)

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = 'Distrito'
        verbose_name_plural = 'Distritos'
        ordering = ['name']
        db_table = 'district'
