<?php

namespace App\Api\Controllers;

use Exception;
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
use App\Validation\ExposesPostValidationHook;
use Illuminate\Contracts\Support\Responsable;
use App\Http\V1\Responses\InternalServerError;
use Symfony\Component\Validator\Validator\ValidatorInterface;

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
        } catch (Throwable $e) {
            return $this->convertExceptionToResponse($e);
        }

        return null;
    }

    protected function convertExceptionToResponse(Throwable $e): Responsable
    {
        switch (true) {
            case $e instanceof ValidationFailedException:
                return ValidationErrors::new($e->errors());
            case $e instanceof InvalidPropertyException:
                // TODO: need to handle these; may be thrown from Value objects, need API translations
                throw $e;
            case $e instanceof NotFoundException:
                return NotFound::default();
            default:
                Log::error($e);
                return InternalServerError::new();
        }
    }
    // Only really needed for queries, the command bus is handling validation for 
    // commands via the ValidateCommand middleware
    protected function validate(Validatable $validatable): ?Responsable
    {
        $validator = app(ValidatorInterface::class);
        $errors = $validatable->validate($validator);

        if ($errors !== null && ! $errors->isEmpty()) {
            return $this->convertExceptionToResponse($validatable->validationException($errors));
        }

        if ($validatable instanceof ExposesPostValidationHook) {
            $validatable->postValidationHook();
        }

        return null;
    }
}