<?php

namespace App\Commands;

use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use App\Buses\DefinesHandler;
use App\Validation\Validatable;
use App\Api\Translation\HttpField;
use App\Api\Translation\FieldPlacement;
use App\Validation\ValidatesByAttributes;
use App\Handlers\DeleteOrganisationHandler;
use Illuminate\Foundation\Bus\Dispatchable;
use App\Api\Requests\PopulatableFromRequest;
use App\Api\Translation\TranslatesFieldNames;
use App\Validation\ExposesPostValidationHook;
use App\Collections\ValidationErrorCollection;
use Symfony\Component\Validator\Constraints as Assert;
use Symfony\Component\Validator\Validator\ValidatorInterface;


// Left the validation limits as hard-coded here so to not couple this to constants
// in the domain which may change. The API spec shouldn't necessarily 
// change with the domain rules, this is why it's versioned.
class DeleteOrganisationCommand implements PopulatableFromRequest, Validatable, ExposesPostValidationHook, DefinesHandler
{
    use TranslatesFieldNames;
    use ValidatesByAttributes;
    use Dispatchable;

    #[Assert\NotBlank]
    #[Assert\Uuid(versions: [Uuid::ASSERTION_TYPE], message: 'This is not a valid UUID v7.')]
    #[HttpField(name: 'organisation_id', in: FieldPlacement::Uri)]
    protected readonly string $rawId;

    protected readonly Uuid $id;
    
    public function __construct(
        protected ValidatorInterface $validator
    ) {}

    public static function getHandler(): string
    {
        return DeleteOrganisationHandler::class;
    }

    public function populate(Request $request)
    {
        $this->rawId = $request->route()->parameter($this->translate('rawId'), '');
    }

    public function validate(): ?ValidationErrorCollection
    {
        return $this->validateSelf();
    }

    public function id(): Uuid
    {
        return $this->id;
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