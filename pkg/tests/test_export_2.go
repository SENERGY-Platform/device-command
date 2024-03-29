/*
 * Copyright 2022 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package tests

var iotExport2 = `{
"/v3/resources/concepts": [
    {
        "base_characteristic_id": "urn:infai:ses:characteristic:f48d7985-7ee7-4119-a791-bc16a953f440",
        "characteristic_ids": [
            "urn:infai:ses:characteristic:f48d7985-7ee7-4119-a791-bc16a953f440"
        ],
        "conversions": null,
        "creator": "dd69ea0d-f553-4336-80f3-7f4567f85c7b",
        "id": "urn:infai:ses:concept:70b40193-0bed-404c-86ec-5ebe9462712e",
        "name": "Segment Cleaning Information",
        "permission_holders": {
            "admin_users": [
                "dd69ea0d-f553-4336-80f3-7f4567f85c7b"
            ],
            "execute_users": [
                "dd69ea0d-f553-4336-80f3-7f4567f85c7b"
            ],
            "read_users": [
                "dd69ea0d-f553-4336-80f3-7f4567f85c7b"
            ],
            "write_users": [
                "dd69ea0d-f553-4336-80f3-7f4567f85c7b"
            ]
        },
        "permissions": {
            "a": true,
            "r": true,
            "w": true,
            "x": true
        },
        "shared": false
    }
],
"/v3/resources/functions": [
    {
        "concept_id": "urn:infai:ses:concept:70b40193-0bed-404c-86ec-5ebe9462712e",
        "creator": "dd69ea0d-f553-4336-80f3-7f4567f85c7b",
        "description": "",
        "display_name": "Start Segment Cleaning",
        "id": "urn:infai:ses:controlling-function:ced44f01-7328-43e3-8db0-ecd12f448758",
        "name": "Set Start Segment Cleaning",
        "permission_holders": {
            "admin_users": [
                "dd69ea0d-f553-4336-80f3-7f4567f85c7b"
            ],
            "execute_users": [
                "dd69ea0d-f553-4336-80f3-7f4567f85c7b"
            ],
            "read_users": [
                "dd69ea0d-f553-4336-80f3-7f4567f85c7b"
            ],
            "write_users": [
                "dd69ea0d-f553-4336-80f3-7f4567f85c7b"
            ]
        },
        "permissions": {
            "a": true,
            "r": true,
            "w": true,
            "x": true
        },
        "rdf_type": "",
        "shared": false
    }
],
"/aspect-nodes/urn:infai:ses:aspect:4f6b9747-1d30-4db1-bbf8-c5904db32771": {
    "id": "urn:infai:ses:aspect:4f6b9747-1d30-4db1-bbf8-c5904db32771",
    "name": "Segment Cleaning",
    "root_id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
    "parent_id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
    "child_ids": [],
    "ancestor_ids": [
        "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307"
    ],
    "descendent_ids": []
},
    "/aspects/urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307": {
        "id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
        "name": "Cleaning",
        "sub_aspects": [
            {
                "id": "urn:infai:ses:aspect:113c86a1-6292-4d6e-a7ad-d58f12b2211c",
                "name": "Zone Cleaning",
                "sub_aspects": null
            },
            {
                "id": "urn:infai:ses:aspect:4f6b9747-1d30-4db1-bbf8-c5904db32771",
                "name": "Segment Cleaning",
                "sub_aspects": null
            }
        ]
    },
    "/characteristics/urn:infai:ses:characteristic:f48d7985-7ee7-4119-a791-bc16a953f440": {
        "allowed_values": null,
        "display_unit": "",
        "id": "urn:infai:ses:characteristic:f48d7985-7ee7-4119-a791-bc16a953f440",
        "name": "Segment Cleaning Information (Valetudo)",
        "sub_characteristics": [
            {
                "allowed_values": null,
                "display_unit": "",
                "id": "urn:infai:ses:characteristic:b0bf0d79-8a23-40d3-a284-2b87e38138be",
                "name": "segment_ids",
                "sub_characteristics": [
                    {
                        "allowed_values": null,
                        "display_unit": "",
                        "id": "urn:infai:ses:characteristic:802c79a9-c96b-4848-848c-bae29fb00375",
                        "name": "*",
                        "sub_characteristics": [],
                        "type": "https://schema.org/Text"
                    }
                ],
                "type": "https://schema.org/ItemList"
            },
            {
                "allowed_values": null,
                "display_unit": "",
                "id": "urn:infai:ses:characteristic:63769613-47bd-4bb9-9624-083669ee6261",
                "max_value": 4,
                "min_value": 1,
                "name": "iterations",
                "sub_characteristics": [],
                "type": "https://schema.org/Integer",
                "value": 1
            },
            {
                "allowed_values": null,
                "display_unit": "",
                "id": "urn:infai:ses:characteristic:aeacd820-8a9e-4a68-a4c6-412537bb8afb",
                "name": "customOrder",
                "sub_characteristics": [],
                "type": "https://schema.org/Boolean",
                "value": true
            }
        ],
        "type": "https://schema.org/StructuredValue"
    },
     "/concepts/urn:infai:ses:concept:70b40193-0bed-404c-86ec-5ebe9462712e": {
        "base_characteristic_id": "urn:infai:ses:characteristic:f48d7985-7ee7-4119-a791-bc16a953f440",
        "characteristic_ids": [
            "urn:infai:ses:characteristic:f48d7985-7ee7-4119-a791-bc16a953f440"
        ],
        "conversions": null,
        "id": "urn:infai:ses:concept:70b40193-0bed-404c-86ec-5ebe9462712e",
        "name": "Segment Cleaning Information"
    },
    "/device-classes/urn:infai:ses:device-class3b83062a-1599-40c5-adf2-f5c04c262f03": {
        "id": "urn:infai:ses:device-class3b83062a-1599-40c5-adf2-f5c04c262f03",
        "image": "",
        "name": "Rotary Valve"
    },
    "/device-classes/urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86": {
        "id": "urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86",
        "image": "https://i.imgur.com/OZOqLcR.png",
        "name": "Lamp"
    },
    "/device-classes/urn:infai:ses:device-class:1785e500-0d45-4ada-b026-75c4a6342f76": {
        "id": "urn:infai:ses:device-class:1785e500-0d45-4ada-b026-75c4a6342f76",
        "image": "",
        "name": "Motor Vehicle"
    },
    "/device-classes/urn:infai:ses:device-class:3a070dc0-f575-4a3f-9faf-85e00f403f57": {
        "id": "urn:infai:ses:device-class:3a070dc0-f575-4a3f-9faf-85e00f403f57",
        "image": "https://i.imgur.com/lyRHWKi.png",
        "name": "Ventilator"
    },
    "/device-classes/urn:infai:ses:device-class:3e138d25-d5ee-4a89-9a83-630f4308941d": {
        "id": "urn:infai:ses:device-class:3e138d25-d5ee-4a89-9a83-630f4308941d",
        "image": "https://i.imgur.com/YHc7cbe.png",
        "name": "Siren"
    },
    "/device-classes/urn:infai:ses:device-class:42c623c9-977b-4c77-ba64-80fdbc5becc0": {
        "id": "urn:infai:ses:device-class:42c623c9-977b-4c77-ba64-80fdbc5becc0",
        "image": "https://i.imgur.com/2bFomGG.png",
        "name": "Door/Window Contact"
    },
    "/device-classes/urn:infai:ses:device-class:5abacb75-3895-41b7-ac68-00542980b60c": {
        "id": "urn:infai:ses:device-class:5abacb75-3895-41b7-ac68-00542980b60c",
        "image": "",
        "name": "Irrigation Control"
    },
    "/device-classes/urn:infai:ses:device-class:767384f1-83d4-4bed-b127-fc89731298f8": {
        "id": "urn:infai:ses:device-class:767384f1-83d4-4bed-b127-fc89731298f8",
        "image": "https://i.imgur.com/BYPe3ZV.png",
        "name": "Button"
    },
    "/device-classes/urn:infai:ses:device-class:78c02d18-d167-45f5-8a2b-0ad6009b6f36": {
        "id": "urn:infai:ses:device-class:78c02d18-d167-45f5-8a2b-0ad6009b6f36",
        "image": "https://i.imgur.com/pnkSAak.png",
        "name": "Nimbus"
    },
    "/device-classes/urn:infai:ses:device-class:79de1bd9-b933-412d-b98e-4cfe19aa3250": {
        "id": "urn:infai:ses:device-class:79de1bd9-b933-412d-b98e-4cfe19aa3250",
        "image": "https://i.imgur.com/CKfqwwN.png",
        "name": "Smart Plug"
    },
    "/device-classes/urn:infai:ses:device-class:82a213bd-18f0-4506-b567-fe16cab41b84": {
        "id": "urn:infai:ses:device-class:82a213bd-18f0-4506-b567-fe16cab41b84",
        "image": "",
        "name": "Vacuum Cleaner"
    },
    "/device-classes/urn:infai:ses:device-class:8bd38ea2-1835-4a1e-ac02-6b3169513fd3": {
        "id": "urn:infai:ses:device-class:8bd38ea2-1835-4a1e-ac02-6b3169513fd3",
        "image": "https://i.imgur.com/YHc7cbe.png",
        "name": "Air Quality Meter"
    },
    "/device-classes/urn:infai:ses:device-class:997937d6-c5f3-4486-b67c-114675038393": {
        "id": "urn:infai:ses:device-class:997937d6-c5f3-4486-b67c-114675038393",
        "image": "https://i.imgur.com/rkfMAXm.png",
        "name": "Thermostat"
    },
    "/device-classes/urn:infai:ses:device-class:a33bdabe-8c89-4f30-85cc-a5042d30ae6e": {
        "id": "urn:infai:ses:device-class:a33bdabe-8c89-4f30-85cc-a5042d30ae6e",
        "image": "https://i.imgur.com/TP3xaH3.png",
        "name": "ECM"
    },
    "/device-classes/urn:infai:ses:device-class:a61a0d6a-e939-4bf0-ae92-1812d9aabe80": {
        "id": "urn:infai:ses:device-class:a61a0d6a-e939-4bf0-ae92-1812d9aabe80",
        "image": "https://i.imgur.com/YHc7cbe.png",
        "name": "Motion Sensor"
    },
    "/device-classes/urn:infai:ses:device-class:aac12bdd-2dc4-4e33-9647-9b16539aab52": {
        "id": "urn:infai:ses:device-class:aac12bdd-2dc4-4e33-9647-9b16539aab52",
        "image": "https://i.imgur.com/pnkSAak.png",
        "name": "Brick"
    },
    "/device-classes/urn:infai:ses:device-class:ac2ea17c-9818-4e60-b011-8880b2fc94ce": {
        "id": "urn:infai:ses:device-class:ac2ea17c-9818-4e60-b011-8880b2fc94ce",
        "image": "https://i.imgur.com/TP3xaH3.png",
        "name": "Water Meter"
    },
    "/device-classes/urn:infai:ses:device-class:d41b8262-706e-4104-9060-ddcd372cea56": {
        "id": "urn:infai:ses:device-class:d41b8262-706e-4104-9060-ddcd372cea56",
        "image": "https://i.imgur.com/ZLbpNVj.png",
        "name": "Photovoltaik"
    },
    "/device-classes/urn:infai:ses:device-class:de1aaf5f-2cf9-49e8-befe-a6135701421f": {
        "id": "urn:infai:ses:device-class:de1aaf5f-2cf9-49e8-befe-a6135701421f",
        "image": "https://i.imgur.com/ZLbpNVj.png",
        "name": "Weather Station"
    },
    "/device-classes/urn:infai:ses:device-class:e8886926-ef3a-456c-a0f2-0373effb5e43": {
        "id": "urn:infai:ses:device-class:e8886926-ef3a-456c-a0f2-0373effb5e43",
        "image": "",
        "name": "Motorized Curtain"
    },
    "/device-classes/urn:infai:ses:device-class:f00e3915-9f09-4aa4-9aaf-1177459c4312": {
        "id": "urn:infai:ses:device-class:f00e3915-9f09-4aa4-9aaf-1177459c4312",
        "image": "https://i.imgur.com/TP3xaH3.png",
        "name": "Electricity Meter"
    },
    "/device-classes/urn:infai:ses:device-class:fea0ad69-3ef2-4332-9afb-ba5885fa7045": {
        "id": "urn:infai:ses:device-class:fea0ad69-3ef2-4332-9afb-ba5885fa7045",
        "image": "https://i.imgur.com/PhCgcaY.png",
        "name": "Wireless MBus Meter"
    },
    "/device-classes/urn:infai:ses:device-class:ff64280a-58e6-4cf9-9a44-e70d3831a79d": {
        "id": "urn:infai:ses:device-class:ff64280a-58e6-4cf9-9a44-e70d3831a79d",
        "image": "https://i.imgur.com/J2vZL6W.png",
        "name": "Multi Sensor"
    },
"/device-types/urn:infai:ses:device-type:ca30a161-0bd4-49b8-86eb-8c48e29eb34e": {
        "attributes": [
            {
                "key": "senergy/local-mqtt",
                "origin": "web-ui",
                "value": "true"
            }
        ],
        "description": "",
        "device_class_id": "urn:infai:ses:device-class:82a213bd-18f0-4506-b567-fe16cab41b84",
        "id": "urn:infai:ses:device-type:ca30a161-0bd4-49b8-86eb-8c48e29eb34e",
        "name": "Valetudo Vacuum Robot",
        "service_groups": [
            {
                "description": "",
                "key": "5ba73b84-b476-4034-8c4a-7b5efbe3aaa9",
                "name": "Auto Empty Dock Manual Trigger"
            },
            {
                "description": "",
                "key": "abc3fdc8-0571-4c7c-beff-2f09fb1ddd61",
                "name": "Basic control"
            },
            {
                "description": "",
                "key": "2e5ba1f5-bbdf-418b-ad0f-72755c7edc22",
                "name": "Consumables monitoring"
            },
            {
                "description": "",
                "key": "f94ce05e-7ff5-4e9c-9930-189234b46b15",
                "name": "Current Statistics"
            },
            {
                "description": "",
                "key": "b5b42a17-bd20-4003-8c6f-c8fbaf07e591",
                "name": "Fan speed control"
            },
            {
                "description": "",
                "key": "87bd4117-e44f-4a68-98ae-62e34c7928f9",
                "name": "Go to Location"
            },
            {
                "description": "",
                "key": "e55a1f07-b333-468b-bd60-4ece9c93a657",
                "name": "Locate"
            },
            {
                "description": "",
                "key": "94fed18a-b1a2-4519-b730-61caf87f4321",
                "name": "Segment cleaning"
            },
            {
                "description": "",
                "key": "33ce0b95-bd7c-4a31-9109-cddab17498ff",
                "name": "Water grade control"
            },
            {
                "description": "",
                "key": "630eca7f-b546-4874-b731-8edfa3960059",
                "name": "Wi-Fi"
            },
            {
                "description": "",
                "key": "67b5daac-eed3-491c-915b-24a05aba9279",
                "name": "Zone Cleaning"
            },
            {
                "description": "",
                "key": "8b3ebe1a-42a2-480b-876d-bda0c0f7a4e7",
                "name": "Map Data"
            },
            {
                "description": "",
                "key": "4eb9d409-0547-4a85-aa7b-ff73c6e503f4",
                "name": "Status"
            }
        ],
        "services": [
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/AttachmentStateAttribute/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:eb78d4d4-e15f-4f6a-af5e-4f4c57a2d14b",
                "inputs": [],
                "interaction": "event",
                "local_id": "dustbin",
                "name": "Get Attachment State Dust Bin",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:0f2601d9-9bd0-4861-90a3-8ee5d6a52d91",
                            "characteristic_id": "urn:infai:ses:characteristic:7dc1bb7e-b256-408a-a6f9-044dc60fdcf5",
                            "function_id": "urn:infai:ses:measuring-function:079ec1a4-5b33-4785-b7c2-60625be63549",
                            "id": "urn:infai:ses:content-variable:160b917e-7216-4875-ac82-0b3bec175f85",
                            "is_void": false,
                            "name": "attached",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Boolean",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:e91a115c-06eb-4c60-8f02-40cbc562020e",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "4eb9d409-0547-4a85-aa7b-ff73c6e503f4"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/AttachmentStateAttribute/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:a5a07a7e-7e0f-4463-abef-2fd1ae0bf7ce",
                "inputs": [],
                "interaction": "event",
                "local_id": "mop",
                "name": "Get Attachment State Mop",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:0f2601d9-9bd0-4861-90a3-8ee5d6a52d91",
                            "characteristic_id": "urn:infai:ses:characteristic:7dc1bb7e-b256-408a-a6f9-044dc60fdcf5",
                            "function_id": "urn:infai:ses:measuring-function:b2ab73cd-1a7d-4e10-aa45-f4a0a10e5a25",
                            "id": "urn:infai:ses:content-variable:d8ad3fcb-f78e-46f0-953b-a32d9e0b5e32",
                            "is_void": false,
                            "name": "attached",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Boolean",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:d571c3fd-e581-4a00-88ec-3fd88a098d45",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "4eb9d409-0547-4a85-aa7b-ff73c6e503f4"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/AttachmentStateAttribute/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:6f7c57c4-6e86-4b60-8d9b-ac1793b9d746",
                "inputs": [],
                "interaction": "event",
                "local_id": "watertank",
                "name": "Get Attachment State Water Tank",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:0f2601d9-9bd0-4861-90a3-8ee5d6a52d91",
                            "characteristic_id": "urn:infai:ses:characteristic:7dc1bb7e-b256-408a-a6f9-044dc60fdcf5",
                            "function_id": "urn:infai:ses:measuring-function:290215fa-6fa7-4acf-8621-0ddba65761e1",
                            "id": "urn:infai:ses:content-variable:7d4c1409-a36a-4420-b89f-6d87e6a2cd97",
                            "is_void": false,
                            "name": "attached",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Boolean",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:f56bf5ae-1956-4bcf-b782-3c3a00622bb7",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "4eb9d409-0547-4a85-aa7b-ff73c6e503f4"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/BatteryStateAttribute/level"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:70644090-d9f8-463f-bad5-8c1efacb096f",
                "inputs": [],
                "interaction": "event",
                "local_id": "70644090-d9f8-463f-bad5-8c1efacb096f",
                "name": "Get Battery Level",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:861227f6-1523-46a7-b8ab-a4e76f0bdd32",
                            "characteristic_id": "urn:infai:ses:characteristic:46f808f4-bb9e-4cc2-bd50-dc33ca74f273",
                            "function_id": "urn:infai:ses:measuring-function:00549f18-88b5-44c7-adb1-f558e8d53d1d",
                            "id": "urn:infai:ses:content-variable:61833d99-2705-4b15-b29e-6323a2bbf39d",
                            "is_void": false,
                            "name": "percentage",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Float",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:4b95d443-0cb9-4500-9932-a4ec0e320e22",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "4eb9d409-0547-4a85-aa7b-ff73c6e503f4"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/BatteryStateAttribute/status"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:db01e5c7-b5b4-40e5-b06b-5ea9c832ed02",
                "inputs": [],
                "interaction": "event",
                "local_id": "db01e5c7-b5b4-40e5-b06b-5ea9c832ed02",
                "name": "Get Battery Status",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:d4625151-ce27-4620-9b7e-93ded78484f8",
                            "characteristic_id": "urn:infai:ses:characteristic:5b5358c4-efd3-48a5-87b5-a8e0a34ce1ab",
                            "function_id": "urn:infai:ses:measuring-function:35a2bf84-ff3e-4cd4-8e99-4d4bd67d6cbd",
                            "id": "urn:infai:ses:content-variable:26767813-5184-4044-9046-c7018a940620",
                            "is_void": false,
                            "name": "status",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:22d6c4e6-e767-49fb-821c-ec93f35ea396",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "4eb9d409-0547-4a85-aa7b-ff73c6e503f4"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/CurrentStatisticsCapability/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:6ae33947-cf59-448d-a338-9d2a5455bc2a",
                "inputs": [],
                "interaction": "event",
                "local_id": "area",
                "name": "Get Current Statistics Area",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
                            "characteristic_id": "urn:infai:ses:characteristic:733d95d9-f7d7-4f2e-9778-14eed5a91824",
                            "function_id": "urn:infai:ses:measuring-function:f4f74bfc-7a58-42cb-855a-e540d566c2fc",
                            "id": "urn:infai:ses:content-variable:0e58c9e5-946f-4b35-b471-9e8ee4e7f374",
                            "is_void": false,
                            "name": "area",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Float",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:95a50d7f-e7ff-4b3b-af61-3912754f443a",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "f94ce05e-7ff5-4e9c-9930-189234b46b15"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/CurrentStatisticsCapability/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:39bfda40-fed6-491e-866c-3729c71d19b5",
                "inputs": [],
                "interaction": "event",
                "local_id": "time",
                "name": "Get Current Statistics Time",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
                            "characteristic_id": "urn:infai:ses:characteristic:9e1024da-3b60-4531-9f29-464addccb13c",
                            "function_id": "urn:infai:ses:measuring-function:77dbfd2c-770a-4c0c-bc5a-c2dcda878107",
                            "id": "urn:infai:ses:content-variable:beac026a-d127-4b2e-9bc6-ea3cc1feaf81",
                            "is_void": false,
                            "name": "runtime",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Float",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:b2704e7f-bc74-48ca-9562-8008601d3594",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "f94ce05e-7ff5-4e9c-9930-189234b46b15"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/StatusStateAttribute/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:49a66f36-f51c-4b16-9077-81c2aae96c7d",
                "inputs": [],
                "interaction": "event",
                "local_id": "error",
                "name": "Get Error Description",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:861227f6-1523-46a7-b8ab-a4e76f0bdd32",
                            "characteristic_id": "urn:infai:ses:characteristic:5b5358c4-efd3-48a5-87b5-a8e0a34ce1ab",
                            "function_id": "urn:infai:ses:measuring-function:5f27e349-cdfb-4a7e-a9fd-28789308f21b",
                            "id": "urn:infai:ses:content-variable:92e62dbc-ad4b-40cb-952b-fb61b4ed5d3a",
                            "is_void": false,
                            "name": "description",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:2a3ece0b-c0b6-4966-bfa5-6b0e9b9fc28d",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "4eb9d409-0547-4a85-aa7b-ff73c6e503f4"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/FanSpeedControlCapability/preset"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:5ff4ee6a-52fb-4073-96b3-27e93b48d238",
                "inputs": [],
                "interaction": "event",
                "local_id": "5ff4ee6a-52fb-4073-96b3-27e93b48d238",
                "name": "Get Fan Speed",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
                            "characteristic_id": "urn:infai:ses:characteristic:2200a112-8d57-48f4-b111-ca60a1c672dc",
                            "function_id": "urn:infai:ses:measuring-function:f6066d39-ed16-4c69-82aa-e18bcf2be2a7",
                            "id": "urn:infai:ses:content-variable:c1f18ed4-1dce-49d4-b4c5-f79edc6ef3f0",
                            "is_void": false,
                            "name": "preset",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:3d0decbb-b9b0-4442-9923-b6a06ea2b582",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "b5b42a17-bd20-4003-8c6f-c8fbaf07e591"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/WifiConfigurationCapability/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:fce8240f-bcea-40b4-adc5-402112101c9e",
                "inputs": [],
                "interaction": "event",
                "local_id": "frequency",
                "name": "Get Frequency",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:861227f6-1523-46a7-b8ab-a4e76f0bdd32",
                            "characteristic_id": "urn:infai:ses:characteristic:b82b8ff5-6d83-4269-bba1-59a130cad7fb",
                            "function_id": "urn:infai:ses:measuring-function:ad05d2d4-8180-4793-a7ee-03e1c0b9c781",
                            "id": "urn:infai:ses:content-variable:3f7acfeb-4abb-46a3-858a-5f977d66c2f8",
                            "is_void": false,
                            "name": "frequency",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:5501687e-32b6-4c78-9400-94553c988345",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "630eca7f-b546-4874-b731-8edfa3960059"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/WifiConfigurationCapability/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:4a7a3847-7d9a-4990-a6db-e7bb8603439d",
                "inputs": [],
                "interaction": "event",
                "local_id": "ips",
                "name": "Get IP Adresses",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:861227f6-1523-46a7-b8ab-a4e76f0bdd32",
                            "characteristic_id": "urn:infai:ses:characteristic:1fc93184-4209-4754-88f1-7d20d3a6fa56",
                            "function_id": "urn:infai:ses:measuring-function:c1bd9824-1b78-44a0-b14d-966ab0acea70",
                            "id": "urn:infai:ses:content-variable:2d25849e-5e24-4c57-81c7-4279ef63f2db",
                            "is_void": false,
                            "name": "ips",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:6668379d-d345-4f37-8436-1a860ae56347",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "630eca7f-b546-4874-b731-8edfa3960059"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/MapData/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:c7804c50-bac8-4427-8c80-09a36ee52098",
                "inputs": [],
                "interaction": "event",
                "local_id": "segments",
                "name": "Get Map Segments",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:861227f6-1523-46a7-b8ab-a4e76f0bdd32",
                            "characteristic_id": "urn:infai:ses:characteristic:5b5358c4-efd3-48a5-87b5-a8e0a34ce1ab",
                            "function_id": "urn:infai:ses:measuring-function:7bba7a1a-5caa-49eb-958b-2cb643e5fcc7",
                            "id": "urn:infai:ses:content-variable:fe97758c-6d91-4987-a7c4-9d043782a311",
                            "is_void": false,
                            "name": "jsonencoded",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:d6e21c9d-7a8e-4c36-875e-e44f01e04a91",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "8b3ebe1a-42a2-480b-876d-bda0c0f7a4e7"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/WifiConfigurationCapability/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:ee03b196-6ba3-445e-b1db-7af3c68d1c32",
                "inputs": [],
                "interaction": "event",
                "local_id": "signal",
                "name": "Get RSSI",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:861227f6-1523-46a7-b8ab-a4e76f0bdd32",
                            "characteristic_id": "urn:infai:ses:characteristic:ff531480-1a59-4dfb-91a9-656d03421d35",
                            "function_id": "urn:infai:ses:measuring-function:9377c366-2896-486e-9b28-cff237509555",
                            "id": "urn:infai:ses:content-variable:38140955-346c-40e4-89a8-3d7f692b5759",
                            "is_void": false,
                            "name": "signal",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Float",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:c63b3b59-1954-4944-b27f-541bf5ef4e9d",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "630eca7f-b546-4874-b731-8edfa3960059"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/ConsumableMonitoringCapability/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:a01c38fb-76e6-4d8e-9f47-21012b03f46b",
                "inputs": [],
                "interaction": "event",
                "local_id": "brush-side_right",
                "name": "Get Remaining Durability Brush Side Right",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:c46eaf2f-2cfa-43c4-9381-eff0e6c77b7c",
                            "characteristic_id": "urn:infai:ses:characteristic:b36eee5d-52f0-4476-a6f7-6dd03b24e0f8",
                            "function_id": "urn:infai:ses:measuring-function:6a6b66f0-4cde-47cc-84ac-8126e4391fd9",
                            "id": "urn:infai:ses:content-variable:245766f3-800f-48f2-9bad-ce78e6534dd0",
                            "is_void": false,
                            "name": "value",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Float",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:a81e2717-1d10-4102-b36d-e79b4f16c90c",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "2e5ba1f5-bbdf-418b-ad0f-72755c7edc22"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/ConsumableMonitoringCapability/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:dbc0170e-3616-431b-92aa-584af8c295ff",
                "inputs": [],
                "interaction": "event",
                "local_id": "filter-main",
                "name": "Get Remaining Durability Filter Main",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:5dedc63a-747a-4654-a615-bad26582431e",
                            "characteristic_id": "urn:infai:ses:characteristic:b36eee5d-52f0-4476-a6f7-6dd03b24e0f8",
                            "function_id": "urn:infai:ses:measuring-function:6a6b66f0-4cde-47cc-84ac-8126e4391fd9",
                            "id": "urn:infai:ses:content-variable:c0b72186-eb68-44ea-92e1-27518996067e",
                            "is_void": false,
                            "name": "value",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Float",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:3dcadade-61a8-4590-9793-94377c386c7d",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "2e5ba1f5-bbdf-418b-ad0f-72755c7edc22"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/ConsumableMonitoringCapability/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:41b9b19e-82fe-4610-87f4-ef5c7b07bc9a",
                "inputs": [],
                "interaction": "event",
                "local_id": "brush-main",
                "name": "Get Remaining Durability Main Brush",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:5fe2556e-a994-4e92-ade2-67f3a8762d3b",
                            "characteristic_id": "urn:infai:ses:characteristic:b36eee5d-52f0-4476-a6f7-6dd03b24e0f8",
                            "function_id": "urn:infai:ses:measuring-function:6a6b66f0-4cde-47cc-84ac-8126e4391fd9",
                            "id": "urn:infai:ses:content-variable:9fefe7be-a8f2-4d68-8ce4-d7ae1ff1144c",
                            "is_void": false,
                            "name": "value",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Float",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:3739108b-cb42-4fbc-a611-591a5cd5853d",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "2e5ba1f5-bbdf-418b-ad0f-72755c7edc22"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/WifiConfigurationCapability/{{.Service}}"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:e41f8391-80ea-4f3a-83aa-e0be9421297b",
                "inputs": [],
                "interaction": "event",
                "local_id": "ssid",
                "name": "Get SSID",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:861227f6-1523-46a7-b8ab-a4e76f0bdd32",
                            "characteristic_id": "urn:infai:ses:characteristic:5b5358c4-efd3-48a5-87b5-a8e0a34ce1ab",
                            "function_id": "urn:infai:ses:measuring-function:5b8fc671-8350-4a67-bb15-750df40cc62d",
                            "id": "urn:infai:ses:content-variable:bd02168a-8dc8-4197-8de7-0f42fb519be7",
                            "is_void": false,
                            "name": "ssid",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:5727b157-c4ea-4234-b0a4-2145e7fa85de",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "630eca7f-b546-4874-b731-8edfa3960059"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/StatusStateAttribute/status"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:2524ec23-3881-41e9-bc84-a18c6fb2f018",
                "inputs": [],
                "interaction": "event",
                "local_id": "2524ec23-3881-41e9-bc84-a18c6fb2f018",
                "name": "Get Vacuum Status",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:861227f6-1523-46a7-b8ab-a4e76f0bdd32",
                            "characteristic_id": "urn:infai:ses:characteristic:5b5358c4-efd3-48a5-87b5-a8e0a34ce1ab",
                            "function_id": "urn:infai:ses:measuring-function:35a2bf84-ff3e-4cd4-8e99-4d4bd67d6cbd",
                            "id": "urn:infai:ses:content-variable:02972036-8f4c-47e3-a1ab-f5807ea909e9",
                            "is_void": false,
                            "name": "status",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:1ed509e7-c94f-4710-aacb-2df15f24feaf",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "4eb9d409-0547-4a85-aa7b-ff73c6e503f4"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/StatusStateAttribute/detail"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:b1dc963b-0c20-4ceb-976e-454989e384cd",
                "inputs": [],
                "interaction": "event",
                "local_id": "b1dc963b-0c20-4ceb-976e-454989e384cd",
                "name": "Get Vacuum Status Detail",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:861227f6-1523-46a7-b8ab-a4e76f0bdd32",
                            "characteristic_id": "urn:infai:ses:characteristic:5b5358c4-efd3-48a5-87b5-a8e0a34ce1ab",
                            "function_id": "urn:infai:ses:measuring-function:35a2bf84-ff3e-4cd4-8e99-4d4bd67d6cbd",
                            "id": "urn:infai:ses:content-variable:900d58f5-9f63-495b-8dd4-dc29c8891f46",
                            "is_void": false,
                            "name": "status",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:e7df01d9-207e-4911-9549-e43032bcd226",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "4eb9d409-0547-4a85-aa7b-ff73c6e503f4"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/event-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/WaterUsageControlCapability/preset"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:8025c9e4-a631-4e48-8f59-143522408519",
                "inputs": [],
                "interaction": "event",
                "local_id": "8025c9e4-a631-4e48-8f59-143522408519",
                "name": "Get Water Grade",
                "outputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
                            "characteristic_id": "urn:infai:ses:characteristic:2200a112-8d57-48f4-b111-ca60a1c672dc",
                            "function_id": "urn:infai:ses:measuring-function:d019faf2-0805-4bf3-b1d0-2ca1bfeae0cc",
                            "id": "urn:infai:ses:content-variable:741c0bf4-2037-4e2e-9862-f134e4722106",
                            "is_void": false,
                            "name": "preset",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:015bfb41-8be2-46e6-a064-dc95a50c3b3e",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "33ce0b95-bd7c-4a31-9109-cddab17498ff"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/cmd-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/AutoEmptyDockManualTriggerCapability/trigger/set"
                    },
                    {
                        "key": "senergy/local-mqtt/resp-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/AutoEmptyDockManualTriggerCapability/trigger"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:b6989c52-796a-43fc-8711-e7ee518786de",
                "inputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:861227f6-1523-46a7-b8ab-a4e76f0bdd32",
                            "characteristic_id": "",
                            "function_id": "urn:infai:ses:controlling-function:5cf7cbd4-11fb-45b3-93de-7197778973af",
                            "id": "urn:infai:ses:content-variable:b0596df0-6d09-4e4f-91fb-d789b37bc831",
                            "is_void": false,
                            "name": "PERFORM",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": "PERFORM"
                        },
                        "id": "urn:infai:ses:content:501c8423-e872-4d0a-9986-63266e9196a0",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "interaction": "request",
                "local_id": "b6989c52-796a-43fc-8711-e7ee518786de",
                "name": "Set Auto Empty Dock Manual Trigger",
                "outputs": [],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "5ba73b84-b476-4034-8c4a-7b5efbe3aaa9"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/resp-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/FanSpeedControlCapability/preset"
                    },
                    {
                        "key": "senergy/local-mqtt/cmd-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/FanSpeedControlCapability/preset/set"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:3beddfeb-7e68-4eaa-800d-bd25f8198800",
                "inputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
                            "characteristic_id": "urn:infai:ses:characteristic:2200a112-8d57-48f4-b111-ca60a1c672dc",
                            "function_id": "urn:infai:ses:controlling-function:deae8f6f-780c-46c9-b99b-04403d8e53f8",
                            "id": "urn:infai:ses:content-variable:3291906d-0bc7-4de0-a44f-5dbb2b75f823",
                            "is_void": false,
                            "name": "preset",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:d269ebee-11cf-44d4-a319-72e5b5295104",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "interaction": "request",
                "local_id": "3beddfeb-7e68-4eaa-800d-bd25f8198800",
                "name": "Set Fan Speed",
                "outputs": [],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "b5b42a17-bd20-4003-8c6f-c8fbaf07e591"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/cmd-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/BasicControlCapability/operation/set"
                    },
                    {
                        "key": "senergy/local-mqtt/resp-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/BasicControlCapability/operation"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:40e95e01-8683-4791-9ccf-837f697452ba",
                "inputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
                            "characteristic_id": "",
                            "function_id": "urn:infai:ses:controlling-function:b103410b-7a04-4dca-8f09-9fa54037b148",
                            "id": "urn:infai:ses:content-variable:c53eacdc-301a-439d-9484-b2b41a8a737c",
                            "is_void": false,
                            "name": "HOME",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": "HOME"
                        },
                        "id": "urn:infai:ses:content:3484cfe7-8442-4800-bacd-bac631dbc007",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "interaction": "request",
                "local_id": "40e95e01-8683-4791-9ccf-837f697452ba",
                "name": "Set Home",
                "outputs": [],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "abc3fdc8-0571-4c7c-beff-2f09fb1ddd61"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/resp-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/LocateCapability/{{.Service}}"
                    },
                    {
                        "key": "senergy/local-mqtt/cmd-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/LocateCapability/{{.Service}}/set"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:9233b6f5-52db-4b1d-b880-5d1311a08c17",
                "inputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:861227f6-1523-46a7-b8ab-a4e76f0bdd32",
                            "characteristic_id": "",
                            "function_id": "urn:infai:ses:controlling-function:51080271-be9b-4997-b5c0-b83ad1c0157b",
                            "id": "urn:infai:ses:content-variable:4928664e-8421-4445-b6a0-c2c35b2f83cd",
                            "is_void": false,
                            "name": "PERFORM",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": "PERFORM"
                        },
                        "id": "urn:infai:ses:content:23bd5f70-aad5-49c4-bed8-376a22e7a801",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "interaction": "request",
                "local_id": "locate",
                "name": "Set Locate",
                "outputs": [],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "e55a1f07-b333-468b-bd60-4ece9c93a657"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/resp-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/GoToLocationCapability/{{.Service}}"
                    },
                    {
                        "key": "senergy/local-mqtt/cmd-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/GoToLocationCapability/{{.Service}}/set"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:3740d968-0ddc-4355-b95c-01867630f139",
                "inputs": [
                    {
                        "content_variable": {
                            "characteristic_id": "",
                            "id": "urn:infai:ses:content-variable:b9489d4e-48bc-49d8-aa93-99944ff1151f",
                            "is_void": false,
                            "name": "root",
                            "serialization_options": null,
                            "sub_content_variables": [
                                {
                                    "characteristic_id": "",
                                    "id": "urn:infai:ses:content-variable:8a451f3d-a717-4cf7-aa0c-dd6164e46e60",
                                    "is_void": false,
                                    "name": "coordinates",
                                    "serialization_options": null,
                                    "sub_content_variables": [
                                        {
                                            "characteristic_id": "",
                                            "id": "urn:infai:ses:content-variable:d77ac6af-dffb-43ea-b9a6-bad85233dc5c",
                                            "is_void": false,
                                            "name": "x",
                                            "serialization_options": null,
                                            "sub_content_variables": null,
                                            "type": "https://schema.org/Integer",
                                            "value": null
                                        },
                                        {
                                            "characteristic_id": "",
                                            "id": "urn:infai:ses:content-variable:1fdc62f3-2fad-40aa-836b-3bc61303b21c",
                                            "is_void": false,
                                            "name": "y",
                                            "serialization_options": null,
                                            "sub_content_variables": null,
                                            "type": "https://schema.org/Integer",
                                            "value": null
                                        }
                                    ],
                                    "type": "https://schema.org/StructuredValue",
                                    "value": null
                                }
                            ],
                            "type": "https://schema.org/StructuredValue",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:246ca394-0a62-44db-806c-6309df3c9df6",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "interaction": "request",
                "local_id": "go",
                "name": "Set Location",
                "outputs": [],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "87bd4117-e44f-4a68-98ae-62e34c7928f9"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/cmd-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/BasicControlCapability/operation/set"
                    },
                    {
                        "key": "senergy/local-mqtt/resp-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/BasicControlCapability/operation"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:6c05fabd-817e-40b7-a6cc-0192f7724472",
                "inputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
                            "characteristic_id": "",
                            "function_id": "urn:infai:ses:controlling-function:6cf43287-a134-4b53-9377-630aea47142b",
                            "id": "urn:infai:ses:content-variable:595cf6b0-90f6-45ad-aca5-1a3e515b0863",
                            "is_void": false,
                            "name": "PAUSE",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": "PAUSE"
                        },
                        "id": "urn:infai:ses:content:e8082f3d-5a36-42c7-b1be-ec975de2d4b5",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "interaction": "request",
                "local_id": "6c05fabd-817e-40b7-a6cc-0192f7724472",
                "name": "Set Pause",
                "outputs": [],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "abc3fdc8-0571-4c7c-beff-2f09fb1ddd61"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/resp-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/MapSegmentationCapability/clean"
                    },
                    {
                        "key": "senergy/local-mqtt/cmd-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/MapSegmentationCapability/clean/set"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:31712f59-1f7e-445e-8bf1-523cdc5c3a96",
                "inputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:4f6b9747-1d30-4db1-bbf8-c5904db32771",
                            "characteristic_id": "urn:infai:ses:characteristic:f48d7985-7ee7-4119-a791-bc16a953f440",
                            "function_id": "urn:infai:ses:controlling-function:ced44f01-7328-43e3-8db0-ecd12f448758",
                            "id": "urn:infai:ses:content-variable:ce52b962-55a6-46e2-bd96-a9301b83e796",
                            "is_void": false,
                            "name": "root",
                            "serialization_options": null,
                            "sub_content_variables": [
                                {
                                    "characteristic_id": "urn:infai:ses:characteristic:b0bf0d79-8a23-40d3-a284-2b87e38138be",
                                    "id": "urn:infai:ses:content-variable:a3e5c12b-f5e9-4d5a-a35f-386ffecab621",
                                    "is_void": false,
                                    "name": "segment_ids",
                                    "serialization_options": null,
                                    "sub_content_variables": [
                                        {
                                            "characteristic_id": "urn:infai:ses:characteristic:802c79a9-c96b-4848-848c-bae29fb00375",
                                            "id": "urn:infai:ses:content-variable:d671767d-3b2c-4e2f-b8b4-b852d4e68c31",
                                            "is_void": false,
                                            "name": "*",
                                            "serialization_options": null,
                                            "sub_content_variables": null,
                                            "type": "https://schema.org/Text",
                                            "value": null
                                        }
                                    ],
                                    "type": "https://schema.org/ItemList",
                                    "value": null
                                },
                                {
                                    "characteristic_id": "urn:infai:ses:characteristic:63769613-47bd-4bb9-9624-083669ee6261",
                                    "id": "urn:infai:ses:content-variable:06288052-2389-42b9-95b0-26a49e41f7a4",
                                    "is_void": false,
                                    "name": "iterations",
                                    "serialization_options": null,
                                    "sub_content_variables": null,
                                    "type": "https://schema.org/Integer",
                                    "value": null
                                },
                                {
                                    "characteristic_id": "urn:infai:ses:characteristic:aeacd820-8a9e-4a68-a4c6-412537bb8afb",
                                    "id": "urn:infai:ses:content-variable:f2819539-c232-49cb-905c-4064f3c59525",
                                    "is_void": false,
                                    "name": "customOrder",
                                    "serialization_options": null,
                                    "sub_content_variables": null,
                                    "type": "https://schema.org/Boolean",
                                    "value": null
                                }
                            ],
                            "type": "https://schema.org/StructuredValue",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:5d2a2c60-58f2-44f3-8eae-1784c70abdf5",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "interaction": "request",
                "local_id": "31712f59-1f7e-445e-8bf1-523cdc5c3a96",
                "name": "Set Segment Cleaning",
                "outputs": [],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "94fed18a-b1a2-4519-b730-61caf87f4321"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/cmd-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/BasicControlCapability/operation/set"
                    },
                    {
                        "key": "senergy/local-mqtt/resp-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/BasicControlCapability/operation"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:c513ac11-c05a-48ef-9784-93f369e50845",
                "inputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
                            "characteristic_id": "",
                            "function_id": "urn:infai:ses:controlling-function:1664d813-e30d-4f57-a314-aa53925f3dc4",
                            "id": "urn:infai:ses:content-variable:c10a3952-61c7-45fc-af14-26e4c7799b19",
                            "is_void": false,
                            "name": "START",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": "START"
                        },
                        "id": "urn:infai:ses:content:90a8c522-3532-4e1d-8074-4dfc5025e32c",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "interaction": "request",
                "local_id": "c513ac11-c05a-48ef-9784-93f369e50845",
                "name": "Set Start",
                "outputs": [],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "abc3fdc8-0571-4c7c-beff-2f09fb1ddd61"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/cmd-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/BasicControlCapability/operation/set"
                    },
                    {
                        "key": "senergy/local-mqtt/resp-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/BasicControlCapability/operation"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:d58c3c51-675c-4fa7-9647-2fecb9113cc4",
                "inputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
                            "characteristic_id": "",
                            "function_id": "urn:infai:ses:controlling-function:0577719f-e0ab-4a24-ac4f-e7d267ac634a",
                            "id": "urn:infai:ses:content-variable:817221ae-269a-4d81-b7e9-245f773f9b44",
                            "is_void": false,
                            "name": "STOP",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": "STOP"
                        },
                        "id": "urn:infai:ses:content:9062e1ef-f8a3-4d3f-86dd-a0da2ab94831",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "interaction": "request",
                "local_id": "d58c3c51-675c-4fa7-9647-2fecb9113cc4",
                "name": "Set Stop",
                "outputs": [],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "abc3fdc8-0571-4c7c-beff-2f09fb1ddd61"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/resp-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/WaterUsageControlCapability/preset"
                    },
                    {
                        "key": "senergy/local-mqtt/cmd-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/WaterUsageControlCapability/preset/set"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:d852f199-9e80-4e10-9873-cc9bddcee582",
                "inputs": [
                    {
                        "content_variable": {
                            "aspect_id": "urn:infai:ses:aspect:1273e874-355b-48af-9e82-1aa5f8a52307",
                            "characteristic_id": "urn:infai:ses:characteristic:2200a112-8d57-48f4-b111-ca60a1c672dc",
                            "function_id": "urn:infai:ses:controlling-function:fdd396a0-92f4-4766-bd6e-d7c68333afaa",
                            "id": "urn:infai:ses:content-variable:fa96bbd4-5ad8-497b-8e7d-4ddd8037e948",
                            "is_void": false,
                            "name": "preset",
                            "serialization_options": null,
                            "sub_content_variables": null,
                            "type": "https://schema.org/Text",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:7185de1e-8736-427b-9fde-86c19e5fb2ca",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "plain-text"
                    }
                ],
                "interaction": "request",
                "local_id": "d852f199-9e80-4e10-9873-cc9bddcee582",
                "name": "Set Water Grade",
                "outputs": [],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "33ce0b95-bd7c-4a31-9109-cddab17498ff"
            },
            {
                "attributes": [
                    {
                        "key": "senergy/local-mqtt/resp-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/ZoneCleaningCapability/start"
                    },
                    {
                        "key": "senergy/local-mqtt/cmd-topic-tmpl",
                        "origin": "web-ui",
                        "value": "{{.Prefix}}{{.Device}}/ZoneCleaningCapability/start/set"
                    }
                ],
                "description": "",
                "id": "urn:infai:ses:service:d1df50a2-9177-40c2-abe1-8c8751cd32a6",
                "inputs": [
                    {
                        "content_variable": {
                            "characteristic_id": "",
                            "id": "urn:infai:ses:content-variable:78e2c498-675e-427d-867c-0b76822ef7ab",
                            "is_void": false,
                            "name": "root",
                            "serialization_options": null,
                            "sub_content_variables": [
                                {
                                    "characteristic_id": "",
                                    "id": "urn:infai:ses:content-variable:d130b2f1-f686-437f-aed7-10ac63c7b820",
                                    "is_void": false,
                                    "name": "zones",
                                    "serialization_options": null,
                                    "sub_content_variables": [
                                        {
                                            "characteristic_id": "",
                                            "id": "urn:infai:ses:content-variable:da3c759b-291d-46a1-a22d-56d562ba97ad",
                                            "is_void": false,
                                            "name": "*",
                                            "serialization_options": null,
                                            "sub_content_variables": [
                                                {
                                                    "characteristic_id": "",
                                                    "id": "urn:infai:ses:content-variable:3c2ee27c-eec0-4b1a-98f3-b83deae67b28",
                                                    "is_void": false,
                                                    "name": "iterations",
                                                    "serialization_options": null,
                                                    "sub_content_variables": null,
                                                    "type": "https://schema.org/Integer",
                                                    "value": null
                                                },
                                                {
                                                    "characteristic_id": "",
                                                    "id": "urn:infai:ses:content-variable:aaa2865c-d35f-4252-a860-e96f1733b7cf",
                                                    "is_void": false,
                                                    "name": "points",
                                                    "serialization_options": null,
                                                    "sub_content_variables": [
                                                        {
                                                            "characteristic_id": "",
                                                            "id": "urn:infai:ses:content-variable:2ea17e6f-b6a0-4d61-9bfe-d04dad3122f2",
                                                            "is_void": false,
                                                            "name": "pA",
                                                            "serialization_options": null,
                                                            "sub_content_variables": [
                                                                {
                                                                    "characteristic_id": "",
                                                                    "id": "urn:infai:ses:content-variable:0d9bb4b5-1605-4574-8fbb-106b5f406a22",
                                                                    "is_void": false,
                                                                    "name": "x",
                                                                    "serialization_options": null,
                                                                    "sub_content_variables": null,
                                                                    "type": "https://schema.org/Integer",
                                                                    "value": null
                                                                },
                                                                {
                                                                    "characteristic_id": "",
                                                                    "id": "urn:infai:ses:content-variable:af9d02fd-db9d-4b57-926f-5c8bb220f91d",
                                                                    "is_void": false,
                                                                    "name": "y",
                                                                    "serialization_options": null,
                                                                    "sub_content_variables": null,
                                                                    "type": "https://schema.org/Integer",
                                                                    "value": null
                                                                }
                                                            ],
                                                            "type": "https://schema.org/StructuredValue",
                                                            "value": null
                                                        },
                                                        {
                                                            "characteristic_id": "",
                                                            "id": "urn:infai:ses:content-variable:e4fb303a-7da8-435a-b3da-be4b55746283",
                                                            "is_void": false,
                                                            "name": "pB",
                                                            "serialization_options": null,
                                                            "sub_content_variables": [
                                                                {
                                                                    "characteristic_id": "",
                                                                    "id": "urn:infai:ses:content-variable:a213f9e3-9d16-41d6-ac07-720f6e2aa6bd",
                                                                    "is_void": false,
                                                                    "name": "x",
                                                                    "serialization_options": null,
                                                                    "sub_content_variables": null,
                                                                    "type": "https://schema.org/Integer",
                                                                    "value": null
                                                                },
                                                                {
                                                                    "characteristic_id": "",
                                                                    "id": "urn:infai:ses:content-variable:d5c0efcb-1c30-4deb-8e29-5027ac393140",
                                                                    "is_void": false,
                                                                    "name": "y",
                                                                    "serialization_options": null,
                                                                    "sub_content_variables": null,
                                                                    "type": "https://schema.org/Integer",
                                                                    "value": null
                                                                }
                                                            ],
                                                            "type": "https://schema.org/StructuredValue",
                                                            "value": null
                                                        },
                                                        {
                                                            "characteristic_id": "",
                                                            "id": "urn:infai:ses:content-variable:29517ff5-e890-43a4-ac87-187bc0f00b9c",
                                                            "is_void": false,
                                                            "name": "pC",
                                                            "serialization_options": null,
                                                            "sub_content_variables": [
                                                                {
                                                                    "characteristic_id": "",
                                                                    "id": "urn:infai:ses:content-variable:a336b2a3-4513-400e-8305-5feb25959481",
                                                                    "is_void": false,
                                                                    "name": "x",
                                                                    "serialization_options": null,
                                                                    "sub_content_variables": null,
                                                                    "type": "https://schema.org/Integer",
                                                                    "value": null
                                                                },
                                                                {
                                                                    "characteristic_id": "",
                                                                    "id": "urn:infai:ses:content-variable:6449fb8b-bbc6-4c96-8710-3bfcf360208f",
                                                                    "is_void": false,
                                                                    "name": "y",
                                                                    "serialization_options": null,
                                                                    "sub_content_variables": null,
                                                                    "type": "https://schema.org/Integer",
                                                                    "value": null
                                                                }
                                                            ],
                                                            "type": "https://schema.org/StructuredValue",
                                                            "value": null
                                                        },
                                                        {
                                                            "characteristic_id": "",
                                                            "id": "urn:infai:ses:content-variable:31001481-770f-4b36-8392-8fadbe6201d2",
                                                            "is_void": false,
                                                            "name": "pD",
                                                            "serialization_options": null,
                                                            "sub_content_variables": [
                                                                {
                                                                    "characteristic_id": "",
                                                                    "id": "urn:infai:ses:content-variable:cd12f30b-626d-4f5b-ac09-b47992995921",
                                                                    "is_void": false,
                                                                    "name": "x",
                                                                    "serialization_options": null,
                                                                    "sub_content_variables": null,
                                                                    "type": "https://schema.org/Integer",
                                                                    "value": null
                                                                },
                                                                {
                                                                    "characteristic_id": "",
                                                                    "id": "urn:infai:ses:content-variable:21270746-af65-44b0-9777-747c982bb435",
                                                                    "is_void": false,
                                                                    "name": "y",
                                                                    "serialization_options": null,
                                                                    "sub_content_variables": null,
                                                                    "type": "https://schema.org/Integer",
                                                                    "value": null
                                                                }
                                                            ],
                                                            "type": "https://schema.org/StructuredValue",
                                                            "value": null
                                                        }
                                                    ],
                                                    "type": "https://schema.org/StructuredValue",
                                                    "value": null
                                                }
                                            ],
                                            "type": "https://schema.org/StructuredValue",
                                            "value": null
                                        }
                                    ],
                                    "type": "https://schema.org/ItemList",
                                    "value": null
                                }
                            ],
                            "type": "https://schema.org/StructuredValue",
                            "value": null
                        },
                        "id": "urn:infai:ses:content:34945855-4dc7-4581-8628-7e5b10b31239",
                        "protocol_segment_id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                        "serialization": "json"
                    }
                ],
                "interaction": "request",
                "local_id": "d1df50a2-9177-40c2-abe1-8c8751cd32a6",
                "name": "Set Zone Cleaning",
                "outputs": [],
                "protocol_id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
                "service_group_key": "67b5daac-eed3-491c-915b-24a05aba9279"
            }
        ]
    },
"/functions/urn:infai:ses:controlling-function:ced44f01-7328-43e3-8db0-ecd12f448758":{
    "id": "urn:infai:ses:controlling-function:ced44f01-7328-43e3-8db0-ecd12f448758",
    "name": "Set Start Segment Cleaning",
    "display_name": "Start Segment Cleaning",
    "description": "",
    "concept_id": "urn:infai:ses:concept:70b40193-0bed-404c-86ec-5ebe9462712e",
    "rdf_type": "https://senergy.infai.org/ontology/ControllingFunction"
},
    "/protocols/urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e": {
        "handler": "moses",
        "id": "urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e",
        "name": "moses",
        "protocol_segments": [
            {
                "id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a",
                "name": "payload"
            }
        ]
    },
    "/protocols/urn:infai:ses:protocol:6b481216-1619-496d-92fc-a7570ec989e8": {
        "handler": "nimbus",
        "id": "urn:infai:ses:protocol:6b481216-1619-496d-92fc-a7570ec989e8",
        "name": "nimbus",
        "protocol_segments": [
            {
                "id": "urn:infai:ses:protocol-segment:943bdb98-0778-4a94-b3eb-173c31339224",
                "name": "payload"
            }
        ]
    },
    "/protocols/urn:infai:ses:protocol:c9a06d44-0cd0-465b-b0d9-560d604057a2": {
        "handler": "mqtt",
        "id": "urn:infai:ses:protocol:c9a06d44-0cd0-465b-b0d9-560d604057a2",
        "name": "mqtt-connector",
        "protocol_segments": [
            {
                "id": "urn:infai:ses:protocol-segment:ffaaf98e-7360-400c-94d4-7775683d38ca",
                "name": "payload"
            },
            {
                "id": "urn:infai:ses:protocol-segment:e07a3534-d67f-4d35-bc14-352fdccf6d8d",
                "name": "topic"
            }
        ]
    },
    "/protocols/urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b": {
        "handler": "connector",
        "id": "urn:infai:ses:protocol:f3a63aeb-187e-4dd9-9ef5-d97a6eb6292b",
        "name": "standard-connector",
        "protocol_segments": [
            {
                "id": "urn:infai:ses:protocol-segment:9956d8b5-46fa-4381-a227-c1df69808997",
                "name": "metadata"
            },
            {
                "id": "urn:infai:ses:protocol-segment:0d211842-cef8-41ec-ab6b-9dbc31bc3a65",
                "name": "data"
            }
        ]
    }
}`
