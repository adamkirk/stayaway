<?php

namespace App\Http\V1\Controllers;

use Throwable;
use App\Buses\DefaultBus;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Http\V1\Responses\NotFound;
use App\Repositories\Organisations;
use Illuminate\Support\Facades\Log;
use App\Http\V1\Responses\NoContent;
use App\Exceptions\NotFoundException;
use App\Queries\GetOrganisationQuery;
use App\Api\Controllers\BaseController;
use App\Queries\ListOrganisationsQuery;
use App\Commands\CreateOrganisationCommand;
use App\Commands\DeleteOrganisationCommand;
use App\Commands\UpdateOrganisationCommand;
use App\Http\V1\Responses\ValidationErrors;
use App\Exceptions\InvalidPropertyException;
use App\Http\V1\Responses\InternalServerError;
use App\Http\V1\Responses\BadRequestWithErrors;
use App\Http\V1\Responses\Organisations\Created;
use App\Http\V1\Responses\Organisations\Updated;
use App\Http\V1\Responses\Organisations\LoadedMany;
use App\Http\V1\Responses\Organisations\LoadedSingle;
use Symfony\Component\Validator\Validator\ValidatorInterface;

/*
Document why actions are protected, essentially so we can use the __call function 
to intercept them and build the relevant request object.
*/
class OrganisationsController extends BaseController
{
    protected function create(CreateOrganisationCommand $cmd, Organisations $repo): Created|InternalServerError|ValidationErrors
    {
        if (($resp = $this->dispatch($cmd)) !== null) {
            return $resp;
        }

        $org = $repo->byId($cmd->id());

        return Created::fromEntity($org);
    }

    protected function update(UpdateOrganisationCommand $cmd, Organisations $repo): Updated|NotFound|InternalServerError|ValidationErrors
    {
        if (($resp = $this->dispatch($cmd)) !== null) {
            return $resp;
        }

        $org = $repo->byId($cmd->id());

        return Updated::fromEntity($org);
    }

    protected function get(GetOrganisationQuery $query, Organisations $repo): LoadedSingle|NotFound|InternalServerError
    {
        if (($resp = $this->validate($query)) !== null) {
            return $resp;
        }

        $org = $repo->byId($query->id);

        if ($org === null) {
            return NotFound::default();
        }

        return LoadedSingle::fromEntity($org);
    }

    protected function delete(DeleteOrganisationCommand $cmd): NoContent|NotFound|InternalServerError
    {
        if (($resp = $this->dispatch($cmd)) !== null) {
            return $resp;
        }

        return NoContent::new();
    }

    protected function list(ListOrganisationsQuery $query, Organisations $repo): LoadedMany|InternalServerError|ValidationErrors
    {
        if (($resp = $this->validate($query)) !== null) {
            return $resp;
        }

        $page = $repo->page($query->page, $query->pageSize, $query->orderBy, $query->orderDirection);

        return LoadedMany::fromPage($page);
    }
}