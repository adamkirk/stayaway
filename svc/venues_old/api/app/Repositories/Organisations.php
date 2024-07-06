<?php

namespace App\Repositories;

use App\ValueObjects\Uuid;
use App\DTO\OrganisationPage;
use App\Entities\Organisation;
use App\Queries\OrderDirection;
use App\Repositories\SaveResult;

interface Organisations
{
    public function save(Organisation $org): SaveResult;

    public function byId(Uuid $id): ?Organisation;

    public function delete(Uuid $id): DeleteResult;

    public function page(int $page, int $pageSize, string $orderBy, OrderDirection $orderDirection): OrganisationPage;
}