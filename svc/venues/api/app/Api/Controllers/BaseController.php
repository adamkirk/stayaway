<?php

namespace App\Api\Controllers;

use App\Buses\DefaultBus;
use App\Validation\Validatable;
use App\Api\Requests\PopulatableFromRequest;
use App\Validation\ExposesPostValidationHook;

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

            if ($arg instanceof Validatable) {
                $errors = $arg->validate();

                if ($errors !== null) {
                    return $arg->invalidResponse($errors);
                }

                if ($arg instanceof ExposesPostValidationHook) {
                    $arg->postValidationHook();
                }
            }
        }

        return call_user_func_array(array($this, $methodName), $args);
    }
}