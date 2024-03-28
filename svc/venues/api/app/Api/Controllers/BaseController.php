<?php

namespace App\Api\Controllers;

use App\Api\Requests\PopulatableFromRequest;
use App\Api\Requests\Validatable;
use App\Http\V1\Responses\ValidationErrors;

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
                    return ValidationErrors::new($errors);
                }
            }
        }

        return call_user_func_array(array($this, $methodName), $args);
    }
}