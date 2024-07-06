<?php

namespace Infra\Database\Repositories;

use App\ValueObjects\Uuid;
use App\DTO\PaginationMeta;
use App\DTO\OrganisationPage;
use App\Entities\Organisation;
use App\Queries\OrderDirection;
use App\Repositories\SaveResult;
use App\Repositories\DeleteResult;
use App\Repositories\Organisations;
use App\Collections\OrganisationCollection;
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

    public function page(int $page, int $pageSize, string $orderBy, OrderDirection $orderDirection): OrganisationPage
    {
        $pageIndex = $page - 1;

        $totalResults = ElOrganisation::count();

        $models = ElOrganisation::skip($pageIndex * $pageSize)
            ->take($pageSize)
            ->orderBy($orderBy, $orderDirection->value)
            ->get();

        $models = array_map(fn(ElOrganisation $org) => $org->toEntity(), $models->all());

        return new OrganisationPage(
            organisations: OrganisationCollection::fromArray($models),
            pagination: new PaginationMeta(
                page: $page,
                pageSize: $pageSize,
                totalPages: ceil($totalResults / $pageSize),
                totalResults: $totalResults,
            ),
        );
    }

    public function delete(Uuid $id): DeleteResult
    {
        $model = ElOrganisation::find($id->toString());

        if ($model === null) {
            return DeleteResult::NotFound;
        }

        $model->delete();

        return DeleteResult::Deleted;
    }
}