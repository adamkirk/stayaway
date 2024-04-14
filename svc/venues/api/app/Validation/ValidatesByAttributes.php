<?php

namespace App\Validation;

use App\Errors\ErrorType;
use App\Errors\ValidationError;
use App\Collections\ValidationErrorCollection;
use App\Http\V1\Responses\ValidationErrors;
use App\Api\Translation\TranslatesFieldNames;
use Illuminate\Contracts\Support\Responsable;
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

    public function invalidResponse(ValidationErrorCollection $errors): Responsable
    {
        return ValidationErrors::new($errors);
    }
}