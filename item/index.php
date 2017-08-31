<?php
require_once __DIR__ . '/vendor/autoload.php';
require_once __DIR__ . '/GPBMetadata/Proto/Item/Item.php';
require_once __DIR__ . '/Item/Item.php';
require_once __DIR__ . '/Item/GetItemRequest.php';
require_once __DIR__ . '/Item/ListItemsRequest.php';
require_once __DIR__ . '/Item/ListItemsResponse.php';


function GetItem($body) {
    $request = new \Item\GetItemRequest();
    $request->mergeFromString($body);
    $response = new \Item\Item();
    $response->setId("abcdef");
    $response->setName("my item name");
    $response->setDescription("a longer description");
    return $response->serializeToString();
}

$body = file_get_contents('php://input');

$response = "";
switch ($_SERVER["REQUEST_URI"]) {
    case "/item.ItemService/GetItem":
        $response = GetItem($body);
        break;
}

header('Content-Type: application/grpc');
print($response);
