<?php

namespace Infra\Database\Repositories;

use App\ValueObjects\Uuid;
use App\Entities\Organisation;
use App\Repositories\SaveResult;
use App\Repositories\Organisations;
use Infra\Database\Models\Eloquent\Organisation as ElOrganisation;

class EloquentOrganisations implements Organisations
{
    public function save(Organisation $org): SaveResult
    {
        $result = SaveResult::DidNothing;

        /** @var ElOrganisation $model */
        $model = ElOrganisation::find($org->id()->toString());

        if ($model === null) {
            $model = ElOrganisation::fromEntity($org);
            $result = SaveResult::Created;
        } else {
            $model->updateFromEntity($org);
            $result = SaveResult::Updated;
        }

        $model->save();

        return $result;
    }

    public function byId(Uuid $id): ?Organisation
    {
        /** @var ElOrganisation $model */
        $model = ElOrganisation::find($id->toString());

        if ($model === null) {
            return null;
        }

        return $model->toEntity();
    }
}