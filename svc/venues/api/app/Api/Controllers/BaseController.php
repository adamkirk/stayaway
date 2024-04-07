<?php

namespace App\Api\Controllers;

use App\Api\Requests\Validatable;
use App\Api\Requests\PopulatableFromRequest;
use App\Api\Requests\ExposesPostValidationHook;

abstract class BaseController
{
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