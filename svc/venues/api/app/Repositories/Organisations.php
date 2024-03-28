<?php

namespace App\Repositories;

use App\ValueObjects\Uuid;
use App\Entities\Organisation;
use App\Repositories\SaveResult;

interface Organisations
{
    public function save(Organisation $org): SaveResult;

    public function byId(Uuid $id): ?Organisation;
}