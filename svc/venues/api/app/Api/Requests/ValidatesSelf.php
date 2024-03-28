<?php

namespace App\Api\Requests;

use App\Errors\ErrorType;
use App\Errors\ValidationError;
use App\Errors\ValidationErrorCollection;
use App\Api\Translation\TranslatesFieldNames;
use Symfony\Component\Validator\Validator\ValidatorInterface;

trait ValidatesSelf
{
    use TranslatesFieldNames;

    abstract protected function getValidator(): ValidatorInterface;

    public function validate(): ?ValidationErrorCollection
    {
        $errors = ValidationErrorCollection::new();
        $list = $this->validator->validate($this);

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
}