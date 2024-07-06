<?php

namespace Tests;

use PHPUnit\Framework\TestCase;

abstract class UnitTestCase extends TestCase
{
    public function __setUp()
    {
        dump("here");
    }
}
