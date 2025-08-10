from django.db import models


class State(models.Model):
    name = models.CharField(max_length=255)
    acronym = models.CharField(max_length=2)

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = 'Estado'
        verbose_name_plural = 'Estados'
        ordering = ['name']
        db_table = 'state'
