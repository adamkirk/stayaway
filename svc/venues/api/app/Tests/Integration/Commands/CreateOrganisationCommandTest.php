<?php

namespace App\Tests\Integration\Commands;

use App\Errors\ErrorType;
use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use Illuminate\Routing\Route;
use App\Entities\Organisation;
use Tests\IntegrationTestCase;
use App\Errors\ValidationError;
use App\Commands\CreateOrganisationCommand;
use PHPUnit\Framework\MockObject\MockObject;
use App\Collections\ValidationErrorCollection;
use PHPUnit\Framework\Attributes\DataProvider;
use Symfony\Component\Validator\Validator\ValidatorInterface;
use App\ValueObjects\Organisation as VO;

class CreateOrganisationCommandTest extends IntegrationTestCase
{
    protected function getValidator(): ValidatorInterface
    {
        return resolve(ValidatorInterface::class);
    }

    protected function getMockRoute(): MockObject
    {
        return $this->getMockBuilder(Route::class)
            ->disableOriginalConstructor()
            ->getMock();
    }

    protected function getMockRequest(MockObject $route): MockObject
    {
        $req = $this->getMockBuilder(Request::class)
            ->disableOriginalConstructor()
            ->getMock();

        $req->expects($this->any())->method('route')->willReturn($route);

        return $req;
    }

    /**
     * The generatedId should be created each time on construct
     */
    public function test_fields_are_populated_if_set()
    {
        $subj = new CreateOrganisationCommand($this->getValidator());

        $route = $this->getMockRoute();
        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return match($key) {
                'name' => 'my name',
                'slug' => 'some-slug',
            };
        });

        $subj->populate($req);
        $this->assertInstanceOf(Uuid::class, $subj->id());
        $this->assertEquals($subj->name(), VO\Name::new('my name'));
        $this->assertEquals($subj->slug(), VO\Slug::new('some-slug'));
    }

    public function test_validate_no_errors()
    {
        $subj = new CreateOrganisationCommand($this->getValidator());

        $route = $this->getMockRoute();
        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return match($key) {
                'name' => 'my name',
                'slug' => 'some-slug',
            };
        });

        $subj->populate($req);
        $this->assertInstanceOf(Uuid::class, $subj->id());
        $this->assertEquals($subj->name(), VO\Name::new('my name'));
        $this->assertEquals($subj->slug(), VO\Slug::new('some-slug'));

        $this->assertNull($subj->validate());
    }

    public function test_validate_min_length()
    {
        $subj = new CreateOrganisationCommand($this->getValidator());

        $route = $this->getMockRoute();
        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return match($key) {
                'name' => str_repeat('x', 2),
                'slug' => str_repeat('x', 1),
            };
        });

        $subj->populate($req);
        $errors = $subj->validate();
       
        $this->assertInstanceOf(ValidationErrorCollection::class, $errors);
        $this->assertCount(3, $errors);
        $this->assertEquals(
            new ValidationError(
                name: "name",
                errorType: ErrorType::ValueNotAllowed,
                message: "The name must be at least 3 characters long",
            ),
            $errors[0],
        );

        // Seems odd that this also happens, but the pattern must contain at least two characters to be valid
        $this->assertEquals(
            new ValidationError(
                name: "slug",
                errorType: ErrorType::ValueNotAllowed,
                message: "The slug must start and end with a number or letter, and may contain letters, numbers and hyphens",
            ),
            $errors[1],
        );

        $this->assertEquals(
            new ValidationError(
                name: "slug",
                errorType: ErrorType::ValueNotAllowed,
                message: "The slug must be at least 2 characters long",
            ),
            $errors[2],
        );
    }

    public function test_validate_max_length()
    {
        $subj = new CreateOrganisationCommand($this->getValidator());

        $route = $this->getMockRoute();
        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return match($key) {
                'name' => str_repeat('x', 256),
                'slug' => str_repeat('x', 256),
            };
        });

        $subj->populate($req);
        $errors = $subj->validate();
       
        $this->assertInstanceOf(ValidationErrorCollection::class, $errors);
        $this->assertCount(2, $errors);
        $this->assertEquals(
            new ValidationError(
                name: "name",
                errorType: ErrorType::ValueNotAllowed,
                message: "The name cannot be longer than 255 characters",
            ),
            $errors[0],
        );
        $this->assertEquals(
            new ValidationError(
                name: "slug",
                errorType: ErrorType::ValueNotAllowed,
                message: "The slug cannot be longer than 255 characters",
            ),
            $errors[1],
        );
    }

    #[DataProvider('invalidSlugPatterns')]
    public function test_invalid_slug_formats(string $slug)
    {
        $subj = new CreateOrganisationCommand($this->getValidator());

        $route = $this->getMockRoute();
        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) use ($slug) {
            return match($key) {
                'name' => str_repeat('x', 30), // this is fine
                'slug' => $slug,
            };
        });

        $subj->populate($req);
        $errors = $subj->validate();
       
        $this->assertInstanceOf(ValidationErrorCollection::class, $errors);
        $this->assertEquals(
            new ValidationError(
                name: "slug",
                errorType: ErrorType::ValueNotAllowed,
                message: "The slug must start and end with a number or letter, and may contain letters, numbers and hyphens",
            ),
            $errors[0],
        );
    }

    public static function invalidSlugPatterns(): array
    {
        return [
            'spaces' => ["this has a space"],
            'all-caps' => ["THIS-IS-CAPS"],
            'one-cap' => ["this-Is-caps"],
            'starts-with-hyphen' => ["-starts-with-hyphen"],
            'ends-with-hyphen' => ["-ends-with-hyphen-"],
            'starts-and-ends-with-hyphen' => ["-hyphens-everywhere-"],
            'underscores' => ["this_has_underscores"],
            'ampersand' => ["this-has-special&characters"],
            'exclamation' => ["this-has-!!-special-characters"],
            'plus-symbol' => ["this-has-special-characters+"],
            'slashes' => ["this-has-special-//characters"],
        ];
    }

    #[DataProvider('validSlugPatterns')]
    public function test_valid_slug_formats(string $slug)
    {
        $subj = new CreateOrganisationCommand($this->getValidator());

        $route = $this->getMockRoute();
        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) use ($slug) {
            return match($key) {
                'name' => str_repeat('x', 30), // this is fine
                'slug' => $slug,
            };
        });

        $subj->populate($req);
        $errors = $subj->validate();
       
        $this->assertNull($errors);
    }

    public static function validSlugPatterns(): array
    {
        return [
            'lowercase' => ["this-is-a-valid-slug"],
            'no-hyphens' => ["thisisaslug"],
        ];
    }

    public function test_name_is_required_slug_is_not()
    {
        $subj = new CreateOrganisationCommand($this->getValidator());

        $route = $this->getMockRoute();
        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return match($key) {
                'name' => $default,
                'slug' => $default,
            };
        });

        $subj->populate($req);
        $errors = $subj->validate();
       
        $this->assertInstanceOf(ValidationErrorCollection::class, $errors);
        $this->assertCount(2, $errors);
        $this->assertEquals(
            new ValidationError(
                name: "name",
                errorType: ErrorType::ValueNotAllowed,
                message: "This value should not be blank.",
            ),
            $errors[0],
        );

        $this->assertEquals(
            new ValidationError(
                name: "name",
                errorType: ErrorType::ValueNotAllowed,
                message: "The name must be at least 3 characters long",
            ),
            $errors[1],
        );
    }

}
