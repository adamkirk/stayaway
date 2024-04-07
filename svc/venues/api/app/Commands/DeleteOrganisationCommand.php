<?php

namespace App\Commands;

use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use App\Api\Requests\Validatable;
use App\Api\Translation\HttpField;
use App\Api\Requests\ValidatesSelf;
use App\Api\Translation\FieldPlacement;
use Illuminate\Foundation\Bus\Dispatchable;
use App\Api\Requests\PopulatableFromRequest;
use App\Api\Translation\TranslatesFieldNames;
use App\Api\Requests\ExposesPostValidationHook;
use Symfony\Component\Validator\Constraints as Assert;
use Symfony\Component\Validator\Validator\ValidatorInterface;


// Left the validation limits as hard-coded here so to not couple this to constants
// in the domain which may change. The API spec shouldn't necessarily 
// change with the domain rules, this is why it's versioned.
class DeleteOrganisationCommand implements PopulatableFromRequest, Validatable, ExposesPostValidationHook
{
    use TranslatesFieldNames;
    use ValidatesSelf;
    use Dispatchable;

    #[Assert\NotBlank]
    #[Assert\Uuid(versions: [Assert\Uuid::V7_MONOTONIC])]
    #[HttpField(name: 'organisation_id', in: FieldPlacement::Uri)]
    public readonly string $rawId;

    public readonly Uuid $id;
    
    public function __construct(
        protected ValidatorInterface $validator
    ) {}

    public function populate(Request $request)
    {
        $this->rawId = $request->route()->parameter($this->translate('rawId'), '');
    }

    public function postValidationHook(): void
    {
        $this->id = Uuid::fromString($this->rawId);
    }

    protected function getValidator(): ValidatorInterface
    {
        return $this->validator;
    }
}