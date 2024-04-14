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

/*
Document why actions are protected, essentially so we can use the __call function 
to intercept them and build the relevant request object.
*/
class OrganisationsController extends BaseController
{
    protected function create(CreateOrganisationCommand $cmd, Organisations $repo): Created|InternalServerError|ValidationErrors
    {
        try {
            $this->bus->handle($cmd);
        } catch (InvalidPropertyException $e) {
            // Validation error needs translating
            // Shouldn't really be possible as validation should already happened
            // But it also happens later, belt & braces
            // TODO: handle this
            throw $e;
        } catch (Throwable $e) {
            Log::error($e);

            return InternalServerError::new();
        }

        $org = $repo->byId($cmd->id());

        return Created::fromEntity($org);
    }

    protected function update(UpdateOrganisationCommand $cmd, Organisations $repo): Updated|NotFound|InternalServerError|ValidationErrors
    {
        try {
            $this->bus->handle($cmd);
        } catch (InvalidPropertyException $e) {
            // Validation error needs translating
            // Shouldn't really be possible as validation should already happened
            // But it also happens later, belt & braces
            // TODO: handle this
            throw $e;
        } catch (NotFoundException $e) {
            return NotFound::default();
        } catch (Throwable $e) {
            Log::error($e);

            return InternalServerError::new();
        }

        $org = $repo->byId($cmd->id());

        return Updated::fromEntity($org);
    }

    protected function get(GetOrganisationQuery $query, Organisations $repo): LoadedSingle|NotFound|InternalServerError
    {
        $org = $repo->byId($query->id);

        if ($org === null) {
            return NotFound::default();
        }

        return LoadedSingle::fromEntity($org);
    }

    protected function delete(DeleteOrganisationCommand $cmd): NoContent|NotFound|InternalServerError
    {
        try {
            $this->bus->handle($cmd);
        } catch (NotFoundException $e) {
            return NotFound::default();
        } catch (Throwable $e) {
            Log::error($e);

            return InternalServerError::new();
        }

        return NoContent::new();
    }

    protected function list(ListOrganisationsQuery $query, Organisations $repo): LoadedMany|InternalServerError|BadRequestWithErrors
    {
        $page = $repo->page($query->page, $query->pageSize, $query->orderBy, $query->orderDirection);

        return LoadedMany::fromPage($page);
    }
}