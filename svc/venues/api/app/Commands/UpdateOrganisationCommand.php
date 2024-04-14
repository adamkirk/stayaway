<?php

namespace App\Commands;

use App\Errors\ErrorType;
use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use App\Buses\DefinesHandler;
use App\Entities\Organisation;
use App\Errors\ValidationError;
use App\Validation\Validatable;
use App\Api\Translation\HttpField;
use App\Api\Translation\FieldPlacement;
use App\ValueObjects\Organisation as VO;
use App\Validation\ValidatesByAttributes;
use App\Handlers\UpdateOrganisationHandler;
use Illuminate\Foundation\Bus\Dispatchable;
use App\Api\Requests\PopulatableFromRequest;
use App\Api\Translation\TranslatesFieldNames;
use App\Validation\ExposesPostValidationHook;
use App\Collections\ValidationErrorCollection;
use Symfony\Component\Validator\Constraints as Assert;
use Symfony\Component\Validator\Validator\ValidatorInterface;

class UpdateOrganisationCommand implements PopulatableFromRequest, Validatable, ExposesPostValidationHook, DefinesHandler
{
    use TranslatesFieldNames;
    use ValidatesByAttributes;
    use Dispatchable;

    #[Assert\NotBlank]
    #[Assert\Uuid(versions: [Uuid::ASSERTION_TYPE])]
    #[HttpField(name: 'organisation_id', in: FieldPlacement::Uri)]
    protected readonly string $rawId;

    #[Assert\Length(
        min: VO\Name::MIN_LENGTH,
        max: VO\Name::MAX_LENGTH,
        minMessage: 'The name must be at least {{ limit }} characters long',
        maxMessage: 'The name cannot be longer than {{ limit }} characters',
    )]
    protected readonly ?string $name;
    
    #[Assert\Regex(
        pattern: VO\Slug::CHARACTER_SET,
        message: "The slug must start and end with a number or letter, and may contain letters, numbers and hyphens",
    )]
    #[Assert\Length(
        min: VO\Slug::MIN_LENGTH,
        max: VO\Slug::MAX_LENGTH,
        minMessage: 'The slug must be at least {{ limit }} characters long',
        maxMessage: 'The slug cannot be longer than {{ limit }} characters',
    )]
    protected readonly ?string $slug;

    protected readonly Uuid $id;

    public function __construct(
        protected ValidatorInterface $validator
    ) {}

    public static function getHandler(): string
    {
        return UpdateOrganisationHandler::class;
    }

    public function populate(Request $request)
    {
        $this->rawId = $request->route()->parameter($this->translate('rawId'), '');
        $this->name = $request->get($this->translate('name'), null);
        $this->slug = $request->get($this->translate('slug'), null);
    }

    public function validate(): ?ValidationErrorCollection
    {
        $errors = $this->validateSelf();

        if ($errors !== null && ! $errors->isEmpty()) {
            return $errors;
        }

        if ($this->name !== null || $this->slug !== null) {
            return null;
        }

        $errors = ValidationErrorCollection::new();

        $errors->add(new ValidationError($this->translate('name'), ErrorType::Required, "one of these fields must be present"));
        $errors->add(new ValidationError($this->translate('slug'), ErrorType::Required, "one of these fields must be present"));

        return $errors;
    }

    public function postValidationHook(): void
    {
        $this->id = Uuid::fromString($this->rawId);
    }

    public function id(): Uuid
    {
        return $this->id;
    }

    public function name(): ?VO\Name
    {
        return $this->name !== null ? VO\Name::new($this->name) : null;
    }


    public function slug(): ?VO\Slug
    {
        return $this->slug !== null ? VO\Slug::new($this->slug) : null;
    }

    protected function getValidator(): ValidatorInterface
    {
        return $this->validator;
    }
}