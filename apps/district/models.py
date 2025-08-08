from django.db import models

# Create your models here.


class District(models.Model):
    name = models.CharField(max_length=255)
    city = models.ForeignKey('city.City', on_delete=models.CASCADE)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = 'Distrito'
        verbose_name_plural = 'Distritos'
        ordering = ['name']
        db_table = 'district'
        unique_together = ['name', 'city']
        constraints = [
            models.UniqueConstraint(
                fields=['name', 'city'], name='unique_district_city')
        ]
