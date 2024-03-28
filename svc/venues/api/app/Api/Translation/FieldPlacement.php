<?php

namespace App\Api\Translation;

enum FieldPlacement: string
{
    case Query = 'query';
    case Body = 'body';
    case Header = 'header';
    case Uri = 'uri';
}