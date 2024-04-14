<?php

namespace App\Commands\Middleware;

use App\Errors\ErrorType;
use App\Errors\ValidationError;
use App\Validation\Validatable;
use League\Tactician\Middleware;
use Symfony\Component\Validator\Validation;
use App\Api\Translation\TranslatesFieldNames;
use App\Exceptions\ValidationFailedException;
use App\Validation\ExposesPostValidationHook;
use App\Collections\ValidationErrorCollection;

class ValidateCommand implements Middleware
{
    use TranslatesFieldNames;

    public function execute($command, callable $next)
    {
        if (! $command instanceof Validatable) {
            return $next($command);
        }
        
        // not sure if we can inject this and use a single instance or if it stores
        // some kind of state...
        $validator = Validation::createValidatorBuilder()
            ->enableAttributeMapping()
            ->getValidator();

        $errors = $command->validate($validator);

        if ($errors !== null && ! $errors->isEmpty()) {
            throw new ValidationFailedException($errors);
        }

        if ($command instanceof ExposesPostValidationHook) {
            $command->postValidationHook();
        }

        return $next($command);
    }
}