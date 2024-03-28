<?php

namespace App\Http\V1\Controllers;

use Throwable;
use Illuminate\Http\Request;
use Illuminate\Http\JsonResponse;
use App\Repositories\Organisations;
use Illuminate\Support\Facades\Log;
use App\Api\Controllers\BaseController;
use App\Commands\CreateOrganisationCommand;
use App\Http\V1\Responses\ValidationErrors;
use App\Exceptions\InvalidPropertyException;
use App\Http\V1\Responses\InternalServerError;
use App\Http\V1\Responses\Organisations\Created;

/*
Document why actions are protected, essentially so we can use the __call function 
to intercept them and build the relevant request object.
*/
class OrganisationsController extends BaseController
{
    protected function create(CreateOrganisationCommand $cmd, Organisations $repo): Created|InternalServerError|ValidationErrors
    {
        try {
            event($cmd);

        } catch (InvalidPropertyException $e) {
            // Validation error needs translating
            // Shouldn't really be possible as validation should already happened
            // But it also happens later, belt & braces
        } catch (Throwable $e) {
            Log::error($e);

            return InternalServerError::new();
        }

        $org = $repo->byId($cmd->generatedId);

        return Created::fromEntity($org);
    }
}