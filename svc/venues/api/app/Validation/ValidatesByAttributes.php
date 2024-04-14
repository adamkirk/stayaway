<?php

namespace App\Validation;

use Exception;
use App\Errors\ErrorType;
use App\Errors\ValidationError;
use App\Http\V1\Responses\ValidationErrors;
use App\Api\Translation\TranslatesFieldNames;
use App\Exceptions\ValidationFailedException;
use Illuminate\Contracts\Support\Responsable;
use App\Collections\ValidationErrorCollection;
use Symfony\Component\Validator\Validator\ValidatorInterface;

trait ValidatesByAttributes
{
    use TranslatesFieldNames;

    public function validateSelf(ValidatorInterface $validator): ?ValidationErrorCollection
    {
        $errors = ValidationErrorCollection::new();
        $list = $validator->validate($this);

        foreach ($list as $error) {
            $errors->add(new ValidationError(
                $this->translate($error->getPropertyPath()),
                ErrorType::ValueNotAllowed,
                $error->getMessage(),
            ));
        }
        
        if (! $errors->isEmpty()) {
            return $errors;
        }

        return null;
    }

    public function validationException(ValidationErrorCollection $errors): Exception
    {
        return new ValidationFailedException($errors);
    }
}