<?php

namespace App\Queries;

use Exception;
use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use App\Validation\Validatable;
use App\Api\Translation\HttpField;
use App\Http\V1\Responses\NotFound;
use App\Exceptions\NotFoundException;
use App\Api\Translation\FieldPlacement;
use App\Validation\ValidatesByAttributes;
use Illuminate\Foundation\Bus\Dispatchable;
use App\Api\Requests\PopulatableFromRequest;
use App\Validation\ExposesPostValidationHook;
use Illuminate\Contracts\Support\Responsable;
use App\Collections\ValidationErrorCollection;
use Symfony\Component\Validator\Constraints as Assert;
use Symfony\Component\Validator\Validator\ValidatorInterface;

class GetOrganisationQuery implements PopulatableFromRequest, Validatable, ExposesPostValidationHook
{
    use ValidatesByAttributes;
    use Dispatchable;

    #[Assert\NotBlank]
    #[Assert\Uuid(
        versions: [Assert\Uuid::V7_MONOTONIC],
        message:'This is not a valid UUID v7.',
    )]
    #[HttpField(name: 'organisation_id', in: FieldPlacement::Uri)]
    public readonly string $rawId;

    public readonly ?Uuid $id;

    public function populate(Request $request)
    {
        $this->rawId = $request->route()->parameter('organisation_id', '');
    }

    public function validate(ValidatorInterface $validator): ?ValidationErrorCollection
    {
        return $this->validateSelf($validator);
    }

    public function postValidationHook(): void
    {
        $this->id = Uuid::fromString($this->rawId);
    }

    public function validationException(ValidationErrorCollection $errors): Exception
    {
        return new NotFoundException;
    }
}