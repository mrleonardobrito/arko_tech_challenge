from django.db import models


class State(models.Model):
    name = models.CharField(max_length=255)
    acronym = models.CharField(max_length=2)
    created_at = models.DateTimeField(auto_now_add=True)
    updated_at = models.DateTimeField(auto_now=True)

    def __str__(self):
        return self.name

    class Meta:
        verbose_name = 'Estado'
        verbose_name_plural = 'Estados'
        ordering = ['name']
        db_table = 'state'
        unique_together = ['name', 'acronym']
        constraints = [
            models.UniqueConstraint(
                fields=['name', 'acronym'], name='unique_state_acronym')
        ]
