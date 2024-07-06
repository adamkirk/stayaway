<?php

namespace App\Tests\Integration\Commands;

use App\Errors\ErrorType;
use App\ValueObjects\Uuid;
use Illuminate\Http\Request;
use Illuminate\Routing\Route;
use App\Entities\Organisation;
use Tests\IntegrationTestCase;
use App\Errors\ValidationError;
use App\Commands\UpdateOrganisationCommand;
use PHPUnit\Framework\MockObject\MockObject;
use App\Collections\ValidationErrorCollection;
use PHPUnit\Framework\Attributes\DataProviderExternal;
use Symfony\Component\Validator\Validator\ValidatorInterface;
use App\Tests\Integration\Commands\CreateOrganisationCommandTest;
use App\ValueObjects\Organisation as VO;

class UpdateOrganisationCommandTest extends IntegrationTestCase
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
        $uuid = '018eb897-4323-76b3-9c55-483ab7f55f43';
        $subj = new UpdateOrganisationCommand();

        $route = $this->getMockRoute();
        $route->expects($this->exactly(1))->method('parameter')->with('organisation_id')->willReturn($uuid);

        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return match($key) {
                'name' => 'my name',
                'slug' => 'some-slug',
            };
        });

        $subj->populate($req);
        $subj->postValidationHook();
        $this->assertEquals($uuid, $subj->id());
        $this->assertInstanceOf(Uuid::class, $subj->id());
        $this->assertEquals($uuid, $subj->id()->toString());
        $this->assertEquals($subj->name(), VO\Name::new('my name'));
        $this->assertEquals($subj->slug(), VO\Slug::new('some-slug'));
    }

    public function test_validate_no_errors()
    {
        $uuid = '018eb897-4323-76b3-9c55-483ab7f55f43';
        $subj = new UpdateOrganisationCommand();

        $route = $this->getMockRoute();
        $route->expects($this->exactly(1))->method('parameter')->with('organisation_id')->willReturn($uuid);

        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return match($key) {
                'name' => 'my name',
                'slug' => 'some-slug',
            };
        });

        $subj->populate($req);
        $this->assertNull($subj->validate($this->getValidator()));
    }

    public function test_name_can_be_set_without_slug()
    {
        $uuid = '018eb897-4323-76b3-9c55-483ab7f55f43';
        $subj = new UpdateOrganisationCommand();

        $route = $this->getMockRoute();
        $route->expects($this->exactly(1))->method('parameter')->with('organisation_id')->willReturn($uuid);

        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return match($key) {
                'name' => 'my name',
                'slug' => $default,
            };
        });

        $subj->populate($req);
        $this->assertNull($subj->validate($this->getValidator()));
    }


    public function test_slug_can_be_set_without_name()
    {
        $uuid = '018eb897-4323-76b3-9c55-483ab7f55f43';
        $subj = new UpdateOrganisationCommand();

        $route = $this->getMockRoute();
        $route->expects($this->exactly(1))->method('parameter')->with('organisation_id')->willReturn($uuid);

        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return match($key) {
                'name' => $default,
                'slug' => 'some-slug',
            };
        });

        $subj->populate($req);
        $this->assertNull($subj->validate($this->getValidator()));
    }

    public function test_validate_at_least_one_field_must_be_set()
    {
        $uuid = '018eb897-4323-76b3-9c55-483ab7f55f43';
        $subj = new UpdateOrganisationCommand();

        $route = $this->getMockRoute();
        $route->expects($this->exactly(1))->method('parameter')->with('organisation_id')->willReturn($uuid);

        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return $default;
        });

        $subj->populate($req);
        $errors = $subj->validate($this->getValidator());
        $this->assertInstanceOf(ValidationErrorCollection::class, $errors);

        $this->assertEquals(
            new ValidationError(
                name: "name",
                errorType: ErrorType::Required,
                message: "one of these fields must be present",
            ),
            $errors[0],
        );

        $this->assertEquals(
            new ValidationError(
                name: "slug",
                errorType: ErrorType::Required,
                message: "one of these fields must be present",
            ),
            $errors[1],
        );
    }
    public function test_validate_min_length()
    {
        $uuid = '018eb897-4323-76b3-9c55-483ab7f55f43';
        $subj = new UpdateOrganisationCommand();

        $route = $this->getMockRoute();
        $route->expects($this->exactly(1))->method('parameter')->with('organisation_id')->willReturn($uuid);

        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return match($key) {
                'name' => str_repeat('x', 2),
                'slug' => str_repeat('x', 1),
            };
        });

        $subj->populate($req);
        $errors = $subj->validate($this->getValidator());
       
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
        $uuid = '018eb897-4323-76b3-9c55-483ab7f55f43';
        $subj = new UpdateOrganisationCommand();

        $route = $this->getMockRoute();
        $route->expects($this->exactly(1))->method('parameter')->with('organisation_id')->willReturn($uuid);

        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) {
            return match($key) {
                'name' => str_repeat('x', 256),
                'slug' => str_repeat('x', 256),
            };
        });

        $subj->populate($req);
        $errors = $subj->validate($this->getValidator());
       
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

    #[DataProviderExternal(CreateOrganisationCommandTest::class, 'invalidSlugPatterns')]
    public function test_invalid_slug_formats(string $slug)
    {
        $uuid = '018eb897-4323-76b3-9c55-483ab7f55f43';
        $subj = new UpdateOrganisationCommand();

        $route = $this->getMockRoute();
        $route->expects($this->exactly(1))->method('parameter')->with('organisation_id')->willReturn($uuid);

        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) use ($slug) {
            return match($key) {
                // Don't care about name for this test
                'name' => $default,
                'slug' => $slug,
            };
        });

        $subj->populate($req);
        $errors = $subj->validate($this->getValidator());
       
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

    #[DataProviderExternal(CreateOrganisationCommandTest::class, 'validSlugPatterns')]
    public function test_valid_slug_formats(string $slug)
    {
        $uuid = '018eb897-4323-76b3-9c55-483ab7f55f43';
        $subj = new UpdateOrganisationCommand();

        $route = $this->getMockRoute();
        $route->expects($this->exactly(1))->method('parameter')->with('organisation_id')->willReturn($uuid);

        $req = $this->getMockRequest($route);
        $req->expects($this->exactly(2))->method('get')->willReturnCallback(function($key, $default) use ($slug) {
            return match($key) {
                // Don't care about name for this test
                'name' => $default,
                'slug' => $slug,
            };
        });

        $subj->populate($req);
        $errors = $subj->validate($this->getValidator());
       
        $this->assertNull($errors);
    }
}
