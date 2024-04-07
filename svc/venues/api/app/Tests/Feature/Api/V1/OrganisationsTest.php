<?php

namespace App\Tests\Feature\Api\V1;

// use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\IntegrationTestCase;

class OrganisationsTest extends IntegrationTestCase
{
    /**
     * Only tests a simple scenario, more thorough validation testing should 
     * happen at integration or unit level.
     */
    public function test_validation_occurs(): void
    {
        $response = $this->post('/api/v1/organisations', [
            "name" => str_repeat("X", 256), //this is too long
            "slug" => "this-should-be-valid",
        ]);

        $response->assertStatus(422);
        $this->assertEquals([
            "errors" => [
                "name" => [
                    "The name cannot be longer than 255 characters"
                ],
            ],
        ], $response->getData(true));
    }

    public function test_organisation_is_created(): void
    {
        $response = $this->post('/api/v1/organisations', [
            "name" => "Some Name", //this is too long
            "slug" => "this-should-be-valid",
        ]);

        $response->assertStatus(201);
        $body = $response->getData(true);
        $this->assertArrayHasKey("data", $body);
        $data = $body['data'];
        $this->assertArrayHasKey("id", $data);
        $this->assertArrayHasKey("name", $data);
        $this->assertArrayHasKey("slug", $data);

        $this->assertEquals("Some Name", $data['name']);
        $this->assertEquals("this-should-be-valid", $data['slug']);
    }
}
