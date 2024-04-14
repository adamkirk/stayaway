<?php

namespace App\Api\Controllers;

use Throwable;
use App\Buses\DefaultBus;
use App\Validation\Validatable;
use App\Http\V1\Responses\NotFound;
use Illuminate\Support\Facades\Log;
use App\Exceptions\NotFoundException;
use App\Http\V1\Responses\ValidationErrors;
use App\Api\Requests\PopulatableFromRequest;
use App\Exceptions\InvalidPropertyException;
use App\Exceptions\ValidationFailedException;
use Illuminate\Contracts\Support\Responsable;
use App\Http\V1\Responses\InternalServerError;

abstract class BaseController
{
    public function __construct(
        protected DefaultBus $bus
    ) {}

    public function __call($methodName, $args)
    {
        foreach ($args as $arg) {
            if ($arg instanceof PopulatableFromRequest) {
                $arg->populate(request());
            }
        }

        return call_user_func_array(array($this, $methodName), $args);
    }

    protected function dispatch($cmd): ?Responsable
    {
        try {
            $this->bus->handle($cmd);
        } catch (ValidationFailedException $e) {
            return ValidationErrors::new($e->errors());
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

        return null;
    }
}