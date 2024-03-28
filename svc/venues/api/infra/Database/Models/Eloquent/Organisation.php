<?php

namespace Infra\Database\Models\Eloquent;

use App\ValueObjects\Uuid;
use InvalidArgumentException;
use Illuminate\Database\Eloquent\Model;
use App\Entities\Organisation as EOrganisation;
use Illuminate\Database\Eloquent\Concerns\HasUuids;
use Illuminate\Database\Eloquent\Factories\HasFactory;

/**
 * @property uuid $id
 * @property string $name
 * @property string $slug
 */
class Organisation extends Model
{
    use HasFactory;
    use HasUuids;

    public static function fromEntity(EOrganisation $entity): self
    {
        $model = new self;
        $model->id = $entity->id()->toString();
        $model->name = $entity->name();
        $model->slug = $entity->slug();

        return $model;
    }

    public function updateFromEntity(EOrganisation $entity): void
    {
        if ($entity->id()->toString() !== $this->id) {
            throw new InvalidArgumentException("Cannot update model from an entity that has a different id!");
        }
        $this->name = $entity->name();
        $this->slug = $entity->slug();
    }

    public function toEntity(): EOrganisation
    {
        return EOrganisation::new(
            Uuid::fromString($this->id),
            $this->name,
            $this->slug,
        );
    }
}
