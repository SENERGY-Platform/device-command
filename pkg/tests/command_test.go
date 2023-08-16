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

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/SENERGY-Platform/device-command/pkg/api"
	"github.com/SENERGY-Platform/device-command/pkg/command"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/cloud"
	"github.com/SENERGY-Platform/device-command/pkg/command/dependencies/impl/mgw"
	"github.com/SENERGY-Platform/device-command/pkg/configuration"
	"github.com/SENERGY-Platform/external-task-worker/lib/com/kafka"
	"github.com/SENERGY-Platform/external-task-worker/lib/devicerepository/model"
	"github.com/SENERGY-Platform/external-task-worker/lib/messages"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestCommandUnscaled(t *testing.T) {
	testCommand("", false)(t)
}

func TestCommandUnscaledCloud(t *testing.T) {
	testCommand("", true)(t)
}

func TestCommandScaled(t *testing.T) {
	testCommand("scalingSuffix", false)(t)
}

func TestGroupCommand_SNRGY_1883(t *testing.T) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config, err := configuration.Load("../../config.json")
	if err != nil {
		t.Error(err)
		return
	}
	config.Debug = true

	config.ServerPort, err = GetFreePort()
	if err != nil {
		t.Error(err)
		return
	}

	config, err = timescaleEnv(config, ctx, wg, map[string]map[string]map[string]interface{}{})
	if err != nil {
		t.Error(err)
		return
	}

	deviceTypeStr := `{
    "id": "urn:infai:ses:device-type:57871169-38a8-40cd-871b-184b99776ca3",
    "name": "Philips Extended Color Light (moses)",
    "description": "Philips Hue Extended Color Light (moses)",
    "service_groups": [],
    "services": [
        {
            "id": "urn:infai:ses:service:163fa8dc-d919-43d0-aee0-4c4b50e406f8",
            "local_id": "getBrightness",
            "name": "Get Brightness",
            "description": "",
            "interaction": "request",
            "protocol_id": "urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e",
            "inputs": [],
            "outputs": [
                {
                    "id": "urn:infai:ses:content:941a6838-84e1-4521-a5d8-c6b83f2c4844",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:871efa99-163d-4a0f-9b36-7cb346bb8565",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:0e8c968a-141e-43a3-9412-d943082648e2",
                                "name": "brightness",
                                "is_void": false,
                                "type": "https://schema.org/Float",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:46f808f4-bb9e-4cc2-bd50-dc33ca74f273",
                                "value": null,
                                "serialization_options": null,
                                "function_id": "urn:infai:ses:measuring-function:c51a6ea5-90c3-4223-9052-6fe4136386cd",
                                "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                            },
                            {
                                "id": "urn:infai:ses:content-variable:c7f55305-62f1-4cfc-826e-d5d81b9ee22a",
                                "name": "status",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:c0353532-a8fb-4553-a00b-418cb8a80a65",
                                "value": null,
                                "serialization_options": null
                            },
                            {
                                "id": "urn:infai:ses:content-variable:6423a99f-fb43-4024-81f6-fc56ecbd5b97",
                                "name": "time",
                                "is_void": false,
                                "type": "https://schema.org/Text",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:6bc41b45-a9f3-4d87-9c51-dd3e11257800",
                                "value": null,
                                "serialization_options": null,
                                "function_id": "urn:infai:ses:measuring-function:3b4e0766-0d67-4658-b249-295902cd3290",
                                "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                            }
                        ],
                        "characteristic_id": "",
                        "value": null,
                        "serialization_options": null
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "attributes": [],
            "service_group_key": ""
        },
        {
            "id": "urn:infai:ses:service:38af42fa-7e89-4b41-826c-2c851eb6a7e8",
            "local_id": "getColor",
            "name": "Get Color",
            "description": "",
            "interaction": "request",
            "protocol_id": "urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e",
            "inputs": [],
            "outputs": [
                {
                    "id": "urn:infai:ses:content:0b4d2ede-1d5f-4b0b-bb38-5d4c3bb5306d",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:64bae8d4-2508-4407-85fc-82b51de00eac",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:eb0a4432-439d-4af5-af73-db476e151c3f",
                                "name": "red",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:dfe6be4a-650c-4411-8d87-062916b48951",
                                "value": null,
                                "serialization_options": null
                            },
                            {
                                "id": "urn:infai:ses:content-variable:a858090f-f4a5-4d19-94a8-81f59bf937f1",
                                "name": "green",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:5ef27837-4aca-43ad-b8f6-4d95cf9ed99e",
                                "value": null,
                                "serialization_options": null
                            },
                            {
                                "id": "urn:infai:ses:content-variable:55b9c48f-78a8-443d-a042-842ad8d38879",
                                "name": "blue",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:590af9ef-3a5e-4edb-abab-177cb1320b17",
                                "value": null,
                                "serialization_options": null
                            },
                            {
                                "id": "urn:infai:ses:content-variable:5056535c-5bcd-4f6e-bab3-2a7e958fdd34",
                                "name": "status",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:c0353532-a8fb-4553-a00b-418cb8a80a65",
                                "value": null,
                                "serialization_options": null
                            },
                            {
                                "id": "urn:infai:ses:content-variable:b7b69993-ae11-4039-9591-ba2e0fde7611",
                                "name": "time",
                                "is_void": false,
                                "type": "https://schema.org/Text",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:6bc41b45-a9f3-4d87-9c51-dd3e11257800",
                                "value": null,
                                "serialization_options": null,
                                "function_id": "urn:infai:ses:measuring-function:3b4e0766-0d67-4658-b249-295902cd3290",
                                "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                            }
                        ],
                        "characteristic_id": "urn:infai:ses:characteristic:5b4eea52-e8e5-4e80-9455-0382f81a1b43",
                        "value": null,
                        "serialization_options": null,
                        "function_id": "urn:infai:ses:measuring-function:bdb6a7c8-4a3d-4fe0-bab3-ce02e09b5869",
                        "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "attributes": [],
            "service_group_key": ""
        },
        {
            "id": "urn:infai:ses:service:7ac752b5-8d25-48b0-aa8a-07dffeb55347",
            "local_id": "getKelvin",
            "name": "Get Color Temperature",
            "description": "",
            "interaction": "request",
            "protocol_id": "urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e",
            "inputs": [],
            "outputs": [
                {
                    "id": "urn:infai:ses:content:854acb4f-2cd4-48ae-b23b-85c32090e9d6",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:40ab57d5-201f-40d8-9c9e-f9b0f0f5f57d",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:a04e23c7-2292-40db-a9d1-6c826e970dd4",
                                "name": "kelvin",
                                "is_void": false,
                                "type": "https://schema.org/Float",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
                                "value": null,
                                "serialization_options": null,
                                "function_id": "urn:infai:ses:measuring-function:fb0f474f-c1d7-4a90-971b-a0b9915968f0",
                                "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                            },
                            {
                                "id": "urn:infai:ses:content-variable:9ea10499-50d5-4dc6-821a-49893e15a849",
                                "name": "status",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:c0353532-a8fb-4553-a00b-418cb8a80a65",
                                "value": null,
                                "serialization_options": null
                            },
                            {
                                "id": "urn:infai:ses:content-variable:80181cb5-bbcb-4c68-8ad8-fc2f4bfd9a7a",
                                "name": "time",
                                "is_void": false,
                                "type": "https://schema.org/Text",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:6bc41b45-a9f3-4d87-9c51-dd3e11257800",
                                "value": null,
                                "serialization_options": null,
                                "function_id": "urn:infai:ses:measuring-function:3b4e0766-0d67-4658-b249-295902cd3290",
                                "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                            }
                        ],
                        "characteristic_id": "",
                        "value": null,
                        "serialization_options": null
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "attributes": [],
            "service_group_key": ""
        },
        {
            "id": "urn:infai:ses:service:4565b9a6-68c9-4d66-9c5e-67e6fc5e989e",
            "local_id": "getPower",
            "name": "Get Power",
            "description": "",
            "interaction": "request",
            "protocol_id": "urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e",
            "inputs": [],
            "outputs": [
                {
                    "id": "urn:infai:ses:content:a0d22df3-1998-43bf-b48a-84e8c95fa1e9",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:b12505dc-aa57-4cae-b6b9-2bcfbf64b8db",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:4ce88010-cb10-4443-8f0e-d7a1beba61c0",
                                "name": "power",
                                "is_void": false,
                                "type": "https://schema.org/Boolean",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:7dc1bb7e-b256-408a-a6f9-044dc60fdcf5",
                                "value": null,
                                "serialization_options": null,
                                "function_id": "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
                                "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                            },
                            {
                                "id": "urn:infai:ses:content-variable:008b09ad-04ca-42be-aadf-49e4f9809ac8",
                                "name": "status",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:c0353532-a8fb-4553-a00b-418cb8a80a65",
                                "value": null,
                                "serialization_options": null
                            },
                            {
                                "id": "urn:infai:ses:content-variable:30a042d7-6cfb-48ab-9a82-138745b2964f",
                                "name": "time",
                                "is_void": false,
                                "type": "https://schema.org/Text",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:6bc41b45-a9f3-4d87-9c51-dd3e11257800",
                                "value": null,
                                "serialization_options": null,
                                "function_id": "urn:infai:ses:measuring-function:3b4e0766-0d67-4658-b249-295902cd3290",
                                "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                            }
                        ],
                        "characteristic_id": "",
                        "value": null,
                        "serialization_options": null
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "attributes": [],
            "service_group_key": ""
        },
        {
            "id": "urn:infai:ses:service:c8751122-ec97-4a84-82e5-fd0341a294e5",
            "local_id": "setBrightness",
            "name": "Set Brightness",
            "description": "",
            "interaction": "request",
            "protocol_id": "urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e",
            "inputs": [
                {
                    "id": "urn:infai:ses:content:32b69aa9-1fcf-4045-b1c6-67e142862cfc",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:acc1b599-08b2-4691-9126-c2a7b0be36f6",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:cbf646e1-9a21-4cac-9ab3-6b89b7ea6fe6",
                                "name": "brightness",
                                "is_void": false,
                                "type": "https://schema.org/Float",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:46f808f4-bb9e-4cc2-bd50-dc33ca74f273",
                                "value": null,
                                "serialization_options": null,
                                "function_id": "urn:infai:ses:controlling-function:6ce74d6d-7acb-40b7-aac7-49daca214e85",
                                "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                            },
                            {
                                "id": "urn:infai:ses:content-variable:462db592-797d-43df-8059-ab4d06272729",
                                "name": "duration",
                                "is_void": false,
                                "type": "https://schema.org/Float",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:9e1024da-3b60-4531-9f29-464addccb13c",
                                "value": null,
                                "serialization_options": null
                            }
                        ],
                        "characteristic_id": "",
                        "value": null,
                        "serialization_options": null
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "outputs": [
                {
                    "id": "urn:infai:ses:content:14b0e0b6-efad-4684-a8cc-fe8871ec2b0b",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:149eed98-a02b-46f7-be11-d04177873cc4",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:74ef2588-f54d-46c7-add8-d9bef536a0ca",
                                "name": "status",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:c0353532-a8fb-4553-a00b-418cb8a80a65",
                                "value": null,
                                "serialization_options": null
                            }
                        ],
                        "characteristic_id": "",
                        "value": null,
                        "serialization_options": null
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "attributes": [],
            "service_group_key": ""
        },
        {
            "id": "urn:infai:ses:service:0f43f156-cd03-4b81-8196-57673e49b15c",
            "local_id": "setColor",
            "name": "Set Color",
            "description": "",
            "interaction": "request",
            "protocol_id": "urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e",
            "inputs": [
                {
                    "id": "urn:infai:ses:content:80111071-fc44-4d06-b443-14dc33685a67",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:05a0f4a4-4cf9-4e17-8e1a-113776856a67",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:1b522486-0764-486f-9ba6-deecd63e0af0",
                                "name": "red",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:dfe6be4a-650c-4411-8d87-062916b48951",
                                "value": null,
                                "serialization_options": null
                            },
                            {
                                "id": "urn:infai:ses:content-variable:bd617a0c-8761-4e43-8548-38d381871f48",
                                "name": "green",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:5ef27837-4aca-43ad-b8f6-4d95cf9ed99e",
                                "value": null,
                                "serialization_options": null
                            },
                            {
                                "id": "urn:infai:ses:content-variable:7112c185-f473-4923-ade4-761ebf2c8c97",
                                "name": "blue",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:590af9ef-3a5e-4edb-abab-177cb1320b17",
                                "value": null,
                                "serialization_options": null
                            },
                            {
                                "id": "urn:infai:ses:content-variable:4df18031-8902-4437-a63c-4ea90b89536f",
                                "name": "duration",
                                "is_void": false,
                                "type": "https://schema.org/Float",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:9e1024da-3b60-4531-9f29-464addccb13c",
                                "value": null,
                                "serialization_options": null
                            }
                        ],
                        "characteristic_id": "urn:infai:ses:characteristic:5b4eea52-e8e5-4e80-9455-0382f81a1b43",
                        "value": null,
                        "serialization_options": null,
                        "function_id": "urn:infai:ses:controlling-function:c54e2a89-1fb8-4ecb-8993-a7b40b355599",
                        "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "outputs": [
                {
                    "id": "urn:infai:ses:content:37c36a9d-0250-4eeb-a568-480ca692be18",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:12b59e29-adad-49d2-8a8b-37134b4f1e41",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:5513ec3a-47c8-4657-b15b-10bc72c8cfc9",
                                "name": "status",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:c0353532-a8fb-4553-a00b-418cb8a80a65",
                                "value": null,
                                "serialization_options": null
                            }
                        ],
                        "characteristic_id": "",
                        "value": null,
                        "serialization_options": null
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "attributes": [],
            "service_group_key": ""
        },
        {
            "id": "urn:infai:ses:service:6add6a31-c5f3-41e4-8224-3004d4407ec1",
            "local_id": "setKelvin",
            "name": "Set Color Temperature",
            "description": "",
            "interaction": "request",
            "protocol_id": "urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e",
            "inputs": [
                {
                    "id": "urn:infai:ses:content:6fd06e48-eeb4-4935-a8f9-8eab8fff5306",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:7e39c04f-203b-49d2-b540-c8c4870efdfd",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:73caed73-3224-4861-b543-ffd70d70eca0",
                                "name": "kelvin",
                                "is_void": false,
                                "type": "https://schema.org/Float",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
                                "value": null,
                                "serialization_options": null,
                                "function_id": "urn:infai:ses:controlling-function:11ede745-afb3-41a6-98fc-6942d0e9cb33",
                                "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                            },
                            {
                                "id": "urn:infai:ses:content-variable:7bde920a-ca2f-4c57-bf0d-dd866fce0589",
                                "name": "duration",
                                "is_void": false,
                                "type": "https://schema.org/Float",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:9e1024da-3b60-4531-9f29-464addccb13c",
                                "value": null,
                                "serialization_options": null
                            }
                        ],
                        "characteristic_id": "",
                        "value": null,
                        "serialization_options": null
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "outputs": [
                {
                    "id": "urn:infai:ses:content:ca782d35-8ae9-413b-bac9-ad73450b1211",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:bd9becc2-a8fe-42b1-9de0-b26cd58e5e20",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:fdde22dd-60c5-4967-904b-ddecb9dab100",
                                "name": "status",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:c0353532-a8fb-4553-a00b-418cb8a80a65",
                                "value": null,
                                "serialization_options": null
                            }
                        ],
                        "characteristic_id": "",
                        "value": null,
                        "serialization_options": null
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "attributes": [],
            "service_group_key": ""
        },
        {
            "id": "urn:infai:ses:service:59dd05fc-cd67-4f66-98de-bbed8257a868",
            "local_id": "setPower",
            "name": "Set Power Off",
            "description": "",
            "interaction": "request",
            "protocol_id": "urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e",
            "inputs": [
                {
                    "id": "urn:infai:ses:content:9f6ac32e-5f26-423b-947a-8a769773239a",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:b3cd09ff-8d5b-4b35-abad-d998bd46a05f",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:18c19fcf-9ded-4393-995f-6ff751c73d83",
                                "name": "power",
                                "is_void": false,
                                "type": "https://schema.org/Boolean",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:7dc1bb7e-b256-408a-a6f9-044dc60fdcf5",
                                "value": false,
                                "serialization_options": null,
                                "function_id": "urn:infai:ses:controlling-function:2f35150b-9df7-4cad-95bc-165fa00219fd",
                                "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                            }
                        ],
                        "characteristic_id": "",
                        "value": null,
                        "serialization_options": null
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "outputs": [
                {
                    "id": "urn:infai:ses:content:77e18f9c-27f7-4755-82b5-9baba3196666",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:cbab719c-5105-49f0-9419-e3733ffae1b9",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:82d51640-9907-4303-b97c-514e8de7c7ad",
                                "name": "status",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:c0353532-a8fb-4553-a00b-418cb8a80a65",
                                "value": null,
                                "serialization_options": null
                            }
                        ],
                        "characteristic_id": "",
                        "value": null,
                        "serialization_options": null
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "attributes": [],
            "service_group_key": ""
        },
        {
            "id": "urn:infai:ses:service:37ec6222-252a-4025-9f0b-18d2598c47c3",
            "local_id": "setPower",
            "name": "Set Power On",
            "description": "",
            "interaction": "request",
            "protocol_id": "urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e",
            "inputs": [
                {
                    "id": "urn:infai:ses:content:afaf90a4-f76b-44e5-8c81-008b4189cdcb",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:fe0fa8e2-304c-4087-83a5-3f01dcb36b64",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:796bd355-0b42-44de-acef-0c79b912a1b2",
                                "name": "power",
                                "is_void": false,
                                "type": "https://schema.org/Boolean",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:7dc1bb7e-b256-408a-a6f9-044dc60fdcf5",
                                "value": true,
                                "serialization_options": null,
                                "function_id": "urn:infai:ses:controlling-function:79e7914b-f303-4a7d-90af-dee70db05fd9",
                                "aspect_id": "urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6"
                            }
                        ],
                        "characteristic_id": "",
                        "value": null,
                        "serialization_options": null
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "outputs": [
                {
                    "id": "urn:infai:ses:content:d00dd025-e7e7-4e15-bb80-5050f74a1f9a",
                    "content_variable": {
                        "id": "urn:infai:ses:content-variable:638b96e7-8b06-4b47-93b4-bbc983133724",
                        "name": "struct",
                        "is_void": false,
                        "type": "https://schema.org/StructuredValue",
                        "sub_content_variables": [
                            {
                                "id": "urn:infai:ses:content-variable:8a8c558c-a168-4bf0-b7a5-5035d4c88d20",
                                "name": "status",
                                "is_void": false,
                                "type": "https://schema.org/Integer",
                                "sub_content_variables": null,
                                "characteristic_id": "urn:infai:ses:characteristic:c0353532-a8fb-4553-a00b-418cb8a80a65",
                                "value": null,
                                "serialization_options": null
                            }
                        ],
                        "characteristic_id": "",
                        "value": null,
                        "serialization_options": null
                    },
                    "serialization": "json",
                    "protocol_segment_id": "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a"
                }
            ],
            "attributes": [],
            "service_group_key": ""
        }
    ],
    "device_class_id": "urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86",
    "attributes": []
}`
	dt := model.DeviceType{}
	err = json.Unmarshal([]byte(deviceTypeStr), &dt)
	if err != nil {
		t.Error(err)
		return
	}

	devices := map[string]map[string]interface{}{
		"testOwner": {
			"/functions/urn:infai:ses:controlling-function:11ede745-afb3-41a6-98fc-6942d0e9cb33": model.Function{
				Id:          "urn:infai:ses:controlling-function:11ede745-afb3-41a6-98fc-6942d0e9cb33",
				Name:        "Set Color Temperature",
				DisplayName: "",
				Description: "",
				ConceptId:   "urn:infai:ses:concept:efb42538-43d3-41fc-9a65-37f0bd81f97a",
				RdfType:     "https://senergy.infai.org/ontology/ControllingFunction",
			},
			"/concepts/urn:infai:ses:concept:efb42538-43d3-41fc-9a65-37f0bd81f97a": model.Concept{
				Id:                   "urn:infai:ses:concept:efb42538-43d3-41fc-9a65-37f0bd81f97a",
				Name:                 "Color Temperature",
				CharacteristicIds:    []string{"urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683"},
				BaseCharacteristicId: "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
			},
			"/characteristics/urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683": model.Characteristic{
				Id:                 "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
				Name:               "Kelvin (Temperature)",
				DisplayUnit:        "K",
				Type:               "https://schema.org/Float",
				MinValue:           nil,
				MaxValue:           nil,
				Value:              nil,
				SubCharacteristics: nil,
			},
			"/device-classes/urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86": model.DeviceClass{
				Id:    "urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86",
				Image: "https://i.imgur.com/OZOqLcR.png",
				Name:  "Lamp",
			},
			"/devices/urn:infai:ses:device:423a2718-dea0-4f69-85a3-626c52de175b": model.Device{
				Id:           "urn:infai:ses:device:423a2718-dea0-4f69-85a3-626c52de175b",
				LocalId:      "618dfabb-c6a8-4d59-a338-ad9d82735ea2",
				Name:         "lampe 1",
				DeviceTypeId: "urn:infai:ses:device-type:57871169-38a8-40cd-871b-184b99776ca3",
			},
			"/devices/urn:infai:ses:device:9cba27cf-323c-4f55-8c57-0910bc3be990": model.Device{
				Id:           "urn:infai:ses:device:9cba27cf-323c-4f55-8c57-0910bc3be990",
				LocalId:      "52f8272c-e1b3-43e1-9954-7d31470392b9",
				Name:         "lampe 2",
				DeviceTypeId: "urn:infai:ses:device-type:57871169-38a8-40cd-871b-184b99776ca3",
			},
			"/protocols/urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e": model.Protocol{
				Id:               "urn:infai:ses:protocol:3b59ea31-da98-45fd-a354-1b9bd06b837e",
				Name:             "moses",
				Handler:          "moses",
				ProtocolSegments: []model.ProtocolSegment{{Id: "urn:infai:ses:protocol-segment:05f1467c-95c8-4a83-a1ed-1c8369fd158a", Name: "payload"}},
			},
			"/device-types/urn:infai:ses:device-type:57871169-38a8-40cd-871b-184b99776ca3": dt,
			"/device-groups/urn:infai:ses:device-group:be0e1dff-6325-4a1a-abbf-9e442c9cfdc2": model.DeviceGroup{
				Id:       "urn:infai:ses:device-group:be0e1dff-6325-4a1a-abbf-9e442c9cfdc2",
				Name:     "Lampen",
				Image:    "",
				Criteria: nil,
				DeviceIds: []string{
					"urn:infai:ses:device:423a2718-dea0-4f69-85a3-626c52de175b",
					"urn:infai:ses:device:9cba27cf-323c-4f55-8c57-0910bc3be990",
				},
				CriteriaShort: []string{
					"urn:infai:ses:controlling-function:2f35150b-9df7-4cad-95bc-165fa00219fd__urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86_request",
					"urn:infai:ses:controlling-function:79e7914b-f303-4a7d-90af-dee70db05fd9__urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86_request",
					"urn:infai:ses:measuring-function:c51a6ea5-90c3-4223-9052-6fe4136386cd_urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6__request",
					"urn:infai:ses:measuring-function:bdb6a7c8-4a3d-4fe0-bab3-ce02e09b5869_urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6__request",
					"urn:infai:ses:controlling-function:c54e2a89-1fb8-4ecb-8993-a7b40b355599__urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86_request",
					"urn:infai:ses:controlling-function:11ede745-afb3-41a6-98fc-6942d0e9cb33__urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86_request",
					"urn:infai:ses:measuring-function:3b4e0766-0d67-4658-b249-295902cd3290_urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6__request",
					"urn:infai:ses:measuring-function:fb0f474f-c1d7-4a90-971b-a0b9915968f0_urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6__request",
					"urn:infai:ses:controlling-function:6ce74d6d-7acb-40b7-aac7-49daca214e85__urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86_request",
					"urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0_urn:infai:ses:aspect:a7470d73-dde3-41fc-92bd-f16bb28f2da6__request",
				},
			},
		},
	}

	config, err = iotEnv(config, ctx, wg, devices)
	if err != nil {
		t.Error(err)
		return
	}

	config, err = kafkaEnv(config, ctx, wg)
	if err != nil {
		t.Error(err)
		return
	}
	time.Sleep(1 * time.Second)

	//connector mock
	maxWait, err := time.ParseDuration(config.KafkaConsumerMaxWait)
	if err != nil {
		t.Error(err)
		return
	}
	producer, err := kafka.PrepareProducerWithConfig(ctx, config.KafkaUrl, kafka.ProducerConfig{
		AsyncFlushFrequency: 100 * time.Millisecond,
		AsyncCompression:    sarama.CompressionNone,
		SyncCompression:     sarama.CompressionNone,
		Sync:                config.Sync,
		SyncIdempotent:      config.SyncIdempotent,
		PartitionNum:        int(config.PartitionNum),
		ReplicationFactor:   int(config.ReplicationFactor),
		AsyncFlushMessages:  int(config.AsyncFlushMessages),
		TopicConfigMap:      config.KafkaTopicConfigs,
	})
	if err != nil {
		t.Log(config.KafkaUrl)
		t.Error(err)
		return
	}
	producer.Log(log.New(os.Stdout, "[KAFKA-PRODUCER] ", 0))

	err = kafka.NewConsumer(ctx, kafka.ConsumerConfig{
		KafkaUrl:       config.KafkaUrl,
		GroupId:        "test-connector-mock",
		Topic:          "moses",
		MinBytes:       int(config.KafkaConsumerMinBytes),
		MaxBytes:       int(config.KafkaConsumerMaxBytes),
		MaxWait:        maxWait,
		TopicConfigMap: config.KafkaTopicConfigs,
	}, func(_ string, msg []byte, time time.Time) error {
		t.Log("MESSAGE:", string(msg))
		message := messages.ProtocolMsg{}
		err := json.Unmarshal(msg, &message)
		if err != nil {
			t.Error(err)
			return nil
		}
		if !strings.Contains(message.Request.Input["payload"], `"kelvin":5641`) {
			t.Error(message.Request.Input["payload"])
		}
		message.Response.Output = map[string]string{}
		resp, err := json.Marshal(message)
		if err != nil {
			t.Error(err)
			return nil
		}
		err = producer.ProduceWithKey(message.Metadata.ResponseTo, "key", string(resp))
		if err != nil {
			t.Error(err)
			return nil
		}
		return nil
	}, func(err error) {
		t.Error(err)
	})

	cmd, err := command.NewWithFactories(ctx, config, cloud.ComFactory, mgw.MarshallerFactory, cloud.IotFactory, mgw.TimescaleFactory)
	if err != nil {
		t.Error(err)
		return
	}
	err = api.Start(ctx, config, cmd)
	if err != nil {
		t.Error(err)
		return
	}

	time.Sleep(5 * time.Second)

	t.Run("command", sendCommand(config, command.CommandMessage{
		FunctionId:    "urn:infai:ses:controlling-function:11ede745-afb3-41a6-98fc-6942d0e9cb33",
		GroupId:       "urn:infai:ses:device-group:be0e1dff-6325-4a1a-abbf-9e442c9cfdc2",
		DeviceClassId: "urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86",
		Input:         5641,
	}, 200, `[null,null]`))

	t.Run("batch", sendCommandBatch(config, command.BatchRequest{
		{
			FunctionId:    "urn:infai:ses:controlling-function:11ede745-afb3-41a6-98fc-6942d0e9cb33",
			GroupId:       "urn:infai:ses:device-group:be0e1dff-6325-4a1a-abbf-9e442c9cfdc2",
			DeviceClassId: "urn:infai:ses:device-class:14e56881-16f9-4120-bb41-270a43070c86",
			Input:         5641,
		},
	}, 200, `[{"status_code":200,"message":[null,null]}]`))

}

func testCommand(scalingSuffix string, cloudTimescale bool) func(t *testing.T) {
	return func(t *testing.T) {
		wg := &sync.WaitGroup{}
		defer wg.Wait()
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		config, err := configuration.Load("../../config.json")
		if err != nil {
			t.Error(err)
			return
		}
		config.TopicSuffixForScaling = scalingSuffix
		config.Debug = true

		config.ServerPort, err = GetFreePort()
		if err != nil {
			t.Error(err)
			return
		}

		if cloudTimescale == true {
			config, err = timescaleCloudEnv(config, ctx, wg, map[string]map[string]map[string]interface{}{
				"color_event": {
					"urn:infai:ses:service:color_event": {
						"struct.hue":        176,
						"struct.saturation": 70,
						"struct.brightness": 65,
						"struct.on":         true,
						"struct.status":     200,
					},
				},
			})
		} else {
			config, err = timescaleEnv(config, ctx, wg, map[string]map[string]map[string]interface{}{
				"color_event": {
					"getStatus": {
						"hue":        176,
						"saturation": 70,
						"brightness": 65,
						"on":         true,
						"status":     200,
					},
				},
			})
		}

		if err != nil {
			t.Error(err)
			return
		}

		devices := map[string]map[string]interface{}{
			"testOwner": {
				"/devices/urn:infai:ses:device:timestamp-test": model.Device{
					Id:           "urn:infai:ses:device:timestamp-test",
					LocalId:      "d1-timestamp",
					Name:         "d1Name-timestamp",
					DeviceTypeId: "urn:infai:ses:device-type:24b294e8-4676-4782-8dc9-a008c0d94770",
				},
				"/devices/urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866": model.Device{
					Id:           "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
					LocalId:      "d1",
					Name:         "d1Name",
					DeviceTypeId: "urn:infai:ses:device-type:755d892f-ec47-40ce-926a-59201328c138",
				},
				"/devices/temperature2": model.Device{
					Id:           "temperature2",
					LocalId:      "d1",
					Name:         "d1Name",
					DeviceTypeId: "urn:infai:ses:device-type:755d892f-ec47-40ce-926a-59201328c138",
				},
				"/devices/temperature3": model.Device{
					Id:           "temperature3",
					LocalId:      "d1",
					Name:         "d1Name",
					DeviceTypeId: "urn:infai:ses:device-type:755d892f-ec47-40ce-926a-59201328c138",
				},
				"/devices/lamp": model.Device{
					Id:           "lamp",
					LocalId:      "lamp",
					Name:         "lamp",
					DeviceTypeId: "urn:infai:ses:device-type:eb4a3337-01a1-4434-9dcc-064b3955eeef",
				},
				"/devices/lamp2": model.Device{
					Id:           "lamp2",
					LocalId:      "lamp2",
					Name:         "lamp2",
					DeviceTypeId: "urn:infai:ses:device-type:eb4a3337-01a1-4434-9dcc-064b3955eeef",
				},
				"/devices/color_event": model.Device{
					Id:           "color_event",
					LocalId:      "color_event",
					Name:         "color_event",
					DeviceTypeId: "urn:infai:ses:device-type:color_event",
				},
				"/device-groups/group_temperature": model.DeviceGroup{
					Id:        "group_temperature",
					DeviceIds: []string{"urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866", "temperature2", "temperature3"},
				},
				"/device-groups/group_color": model.DeviceGroup{
					Id:        "group_color",
					DeviceIds: []string{"color_event", "lamp", "lamp2"},
				},
			},
		}

		config, err = iotEnv(config, ctx, wg, devices)
		if err != nil {
			t.Error(err)
			return
		}

		config, err = kafkaEnv(config, ctx, wg)
		if err != nil {
			t.Error(err)
			return
		}
		time.Sleep(1 * time.Second)

		//connector mock
		maxWait, err := time.ParseDuration(config.KafkaConsumerMaxWait)
		if err != nil {
			t.Error(err)
			return
		}
		producer, err := kafka.PrepareProducerWithConfig(ctx, config.KafkaUrl, kafka.ProducerConfig{
			AsyncFlushFrequency: 100 * time.Millisecond,
			AsyncCompression:    sarama.CompressionNone,
			SyncCompression:     sarama.CompressionNone,
			Sync:                config.Sync,
			SyncIdempotent:      config.SyncIdempotent,
			PartitionNum:        int(config.PartitionNum),
			ReplicationFactor:   int(config.ReplicationFactor),
			AsyncFlushMessages:  int(config.AsyncFlushMessages),
			TopicConfigMap:      config.KafkaTopicConfigs,
		})
		if err != nil {
			t.Log(config.KafkaUrl)
			t.Error(err)
			return
		}
		producer.Log(log.New(os.Stdout, "[KAFKA-PRODUCER] ", 0))

		serviceCallCount := map[string]int{}

		err = kafka.NewConsumer(ctx, kafka.ConsumerConfig{
			KafkaUrl:       config.KafkaUrl,
			GroupId:        "test-connector-mock",
			Topic:          "connector",
			MinBytes:       int(config.KafkaConsumerMinBytes),
			MaxBytes:       int(config.KafkaConsumerMaxBytes),
			MaxWait:        maxWait,
			TopicConfigMap: config.KafkaTopicConfigs,
		}, func(_ string, msg []byte, time time.Time) error {
			t.Log(string(msg))
			message := messages.ProtocolMsg{}
			err := json.Unmarshal(msg, &message)
			if err != nil {
				t.Error(err)
				return nil
			}
			switch message.Metadata.Service.Id {
			case "urn:infai:ses:service:ec456e2a-81ed-4466-a119-daecfbb2d033":
				message.Response.Output = map[string]string{"data": `{"value": "clear", "lastUpdate": 0, "lastUpdate_unit": "unit"}`}
			case "urn:infai:ses:service:4932d451-3300-4a22-a508-ec740e5789b3":
				if message.Request.Input["data"] != "21" {
					t.Error(message.Request.Input)
				}
			case "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1":
				serviceCallCount[message.Metadata.Service.Id] = serviceCallCount[message.Metadata.Service.Id] + 1
				message.Response.Output = map[string]string{"data": `{"value": 13, "lastUpdate": 42}`}
			case "urn:infai:ses:service:4b6c4567-f256-4dbd-a562-d13442ad4530":
				//create timeout
				return nil
			case "urn:infai:ses:service:36fd778e-b04d-4d72-bed5-1b77ed1164b9":
				//create timeout
				return nil
			case "urn:infai:ses:service:1199edee-fbf7-44fb-9228-9a9db69bbdd4":
				message.Response.Output = map[string]string{"data": `{"brightness": 65, "hue": 176, "saturation": 70, "kelvin": 0, "on": true, "status": 200}`}
			case "urn:infai:ses:service:1b0ef253-16f7-4b65-8a15-fe79fccf7e70":
				input := map[string]float64{}
				err = json.Unmarshal([]byte(message.Request.Input["data"]), &input)
				if err != nil {
					t.Error(err)
					return nil
				}
				//values are truncated to integers, not rounded
				if input["brightness"] != float64(65) {
					t.Error(message.Request.Input)
				}
				if input["hue"] != float64(176) {
					t.Error(message.Request.Input)
				}
				if input["saturation"] != float64(70) {
					t.Error(message.Request.Input)
				}
				if input["duration"] != float64(1) {
					t.Error(message.Request.Input)
				}
			default:
				t.Error("unknown service-id", message.Metadata.Service.Id)
				return nil
			}
			resp, err := json.Marshal(message)
			if err != nil {
				t.Error(err)
				return nil
			}
			err = producer.ProduceWithKey(message.Metadata.ResponseTo, "key", string(resp))
			if err != nil {
				t.Error(err)
				return nil
			}
			return nil
		}, func(err error) {
			t.Error(err)
		})

		timescaleFactory := mgw.TimescaleFactory
		if cloudTimescale {
			timescaleFactory = cloud.TimescaleFactory
		}
		cmd, err := command.NewWithFactories(ctx, config, cloud.ComFactory, mgw.MarshallerFactory, cloud.IotFactory, timescaleFactory)
		if err != nil {
			t.Error(err)
			return
		}
		err = api.Start(ctx, config, cmd)
		if err != nil {
			t.Error(err)
			return
		}

		time.Sleep(5 * time.Second)

		t.Run("device default setTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:controlling-function:99240d90-02dd-4d4f-a47c-069cfe77629c",
			Input:      nil,
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:4932d451-3300-4a22-a508-ec740e5789b3",
		}, 200, "[null]"))

		t.Run("device setTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:controlling-function:99240d90-02dd-4d4f-a47c-069cfe77629c",
			Input:      21,
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:4932d451-3300-4a22-a508-ec740e5789b3",
		}, 200, "[null]"))

		t.Run("device setTemperature in Kelvin", sendCommand(config, command.CommandMessage{
			FunctionId:       "urn:infai:ses:controlling-function:99240d90-02dd-4d4f-a47c-069cfe77629c",
			Input:            294.15,
			DeviceId:         "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:        "urn:infai:ses:service:4932d451-3300-4a22-a508-ec740e5789b3",
			CharacteristicId: "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
		}, 200, "[null]"))

		t.Run("device getTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
		}, 200, "[13]"))

		t.Run("device getTemperature in Kelvin", sendCommand(config, command.CommandMessage{
			FunctionId:       "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			DeviceId:         "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:        "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
			CharacteristicId: "urn:infai:ses:characteristic:75b2d113-1d03-4ef8-977a-8dbcbb31a683",
		}, 200, "[286.15]"))

		t.Run("device timeout", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:00549f18-88b5-44c7-adb1-f558e8d53d1d",
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:36fd778e-b04d-4d72-bed5-1b77ed1164b9",
		}, http.StatusRequestTimeout, `"timeout"`))

		t.Run("invalid command", sendCommand(config, command.CommandMessage{
			FunctionId: "foobar",
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
		}, 500, `"unable to load function: not found"`))

		t.Run("device color", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:controlling-function:c54e2a89-1fb8-4ecb-8993-a7b40b355599",
			Input: map[string]interface{}{
				"r": 50,
				"g": 168,
				"b": 162,
			},
			DeviceId:  "lamp",
			ServiceId: "urn:infai:ses:service:1b0ef253-16f7-4b65-8a15-fe79fccf7e70",
		}, 200, "[null]"))

		t.Run("device event color", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:bdb6a7c8-4a3d-4fe0-bab3-ce02e09b5869",
			DeviceId:   "color_event",
			ServiceId:  "urn:infai:ses:service:color_event",
		}, 200, `[{"b":158,"g":166,"r":50}]`))

		t.Run("device event on", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
			DeviceId:   "color_event",
			ServiceId:  "urn:infai:ses:service:color_event",
		}, 200, `[true]`))

		//some services are called as event (timescale call), some as request
		t.Run("device group color", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:bdb6a7c8-4a3d-4fe0-bab3-ce02e09b5869",
			GroupId:    "group_color",
		}, 200, `[{"b":158,"g":166,"r":50},{"b":158,"g":166,"r":50},{"b":158,"g":166,"r":50}]`))

		//some services return a timeout, some return 13
		t.Run("device group getTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			GroupId:    "group_temperature",
		}, 200, "[13,13,13]"))

		t.Run("device batch", sendCommandBatch(config, command.BatchRequest{
			{
				FunctionId: "urn:infai:ses:controlling-function:99240d90-02dd-4d4f-a47c-069cfe77629c",
				Input:      21,
				DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
				ServiceId:  "urn:infai:ses:service:4932d451-3300-4a22-a508-ec740e5789b3",
			},
			{
				FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
				DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
				ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
			},
			{
				FunctionId: "urn:infai:ses:measuring-function:00549f18-88b5-44c7-adb1-f558e8d53d1d",
				DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
				ServiceId:  "urn:infai:ses:service:36fd778e-b04d-4d72-bed5-1b77ed1164b9",
			},
			{
				FunctionId: "foobar",
				DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
				ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
			},
			{
				FunctionId: "urn:infai:ses:controlling-function:c54e2a89-1fb8-4ecb-8993-a7b40b355599",
				Input: map[string]interface{}{
					"r": 50,
					"g": 168,
					"b": 162,
				},
				DeviceId:  "lamp",
				ServiceId: "urn:infai:ses:service:1b0ef253-16f7-4b65-8a15-fe79fccf7e70",
			},
			{
				FunctionId: "urn:infai:ses:measuring-function:bdb6a7c8-4a3d-4fe0-bab3-ce02e09b5869",
				DeviceId:   "color_event",
				ServiceId:  "urn:infai:ses:service:color_event",
			},
			{
				FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
				DeviceId:   "color_event",
				ServiceId:  "urn:infai:ses:service:color_event",
			},
			{
				FunctionId: "urn:infai:ses:controlling-function:99240d90-02dd-4d4f-a47c-069cfe77629c",
				Input:      21,
				DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
				ServiceId:  "urn:infai:ses:service:4932d451-3300-4a22-a508-ec740e5789b3",
			},
			{
				FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
				DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
				ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
			},
			{
				FunctionId: "urn:infai:ses:measuring-function:20d3c1d3-77d7-4181-a9f3-b487add58cd0",
				DeviceId:   "color_event",
				ServiceId:  "urn:infai:ses:service:color_event",
			},
		}, 200, `[{"status_code":200,"message":[null]},{"status_code":200,"message":[13]},{"status_code":408,"message":"timeout"},{"status_code":500,"message":"unable to load function: not found"},{"status_code":200,"message":[null]},{"status_code":200,"message":[{"b":158,"g":166,"r":50}]},{"status_code":200,"message":[true]},{"status_code":200,"message":[null]},{"status_code":200,"message":[13]},{"status_code":200,"message":[true]}]`))

		t.Run("new timestamp", sendCommandBatch(config, command.BatchRequest{
			{
				FunctionId: "urn:infai:ses:measuring-function:3b4e0766-0d67-4658-b249-295902cd3290",
				DeviceId:   "urn:infai:ses:device:timestamp-test",
				ServiceId:  "urn:infai:ses:service:ec456e2a-81ed-4466-a119-daecfbb2d033",
			},
		}, 200, `[{"status_code":200,"message":["1970-01-01T01:00:00+01:00"]}]`))

		t.Run("check callcount", func(t *testing.T) {
			if serviceCallCount["urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1"] != 6 {
				t.Error(serviceCallCount)
			}
		})

		t.Run("device group air getTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			GroupId:    "group_temperature",
			AspectId:   "urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6",
		}, 200, "[13,13,13]"))

		t.Run("device group outside air getTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			GroupId:    "group_temperature",
			AspectId:   "urn:infai:ses:aspect:outside_air",
		}, 200, "[13,13,13]"))

		t.Run("device group outside foo-aspect getTemperature", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			GroupId:    "group_temperature",
			AspectId:   "urn:infai:ses:aspect:foo-aspect",
		}, 200, "[]"))

		t.Run("device getTemperature with aspect", sendCommand(config, command.CommandMessage{
			FunctionId: "urn:infai:ses:measuring-function:f2769eb9-b6ad-4f7e-bd28-e4ea043d2f8b",
			DeviceId:   "urn:infai:ses:device:a486084b-3323-4cbc-9f6b-d797373ae866",
			ServiceId:  "urn:infai:ses:service:6d6067a3-ed4e-45ec-a7eb-b1695340d2f1",
			AspectId:   "urn:infai:ses:aspect:a14c5efb-b0b6-46c3-982e-9fded75b5ab6",
		}, 200, "[13]"))
	}
}

func sendCommandBatch(config configuration.Config, commandMessage command.BatchRequest, expectedCode int, expectedContent string) func(t *testing.T) {
	return func(t *testing.T) {
		buff := &bytes.Buffer{}
		err := json.NewEncoder(buff).Encode(commandMessage)
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest("POST", "http://localhost:"+config.ServerPort+"/commands/batch?timeout=5s", buff)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", testToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()
		actualContent, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		if resp.StatusCode != expectedCode {
			t.Error(resp.StatusCode, string(actualContent))
			return
		}
		if strings.TrimSpace(string(actualContent)) != expectedContent {
			t.Error("\n" + expectedContent + "\n" + string(actualContent))
		}
	}
}

func sendCommand(config configuration.Config, commandMessage command.CommandMessage, expectedCode int, expectedContent string) func(t *testing.T) {
	return func(t *testing.T) {
		buff := &bytes.Buffer{}
		err := json.NewEncoder(buff).Encode(commandMessage)
		if err != nil {
			t.Error(err)
			return
		}
		req, err := http.NewRequest("POST", "http://localhost:"+config.ServerPort+"/commands?timeout=5s", buff)
		if err != nil {
			t.Error(err)
			return
		}
		req.Header.Set("Authorization", testToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Error(err)
			return
		}
		defer resp.Body.Close()
		actualContent, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
			return
		}
		if resp.StatusCode != expectedCode {
			t.Error(resp.StatusCode, string(actualContent))
			return
		}
		if strings.TrimSpace(string(actualContent)) != expectedContent {
			t.Error("\n" + expectedContent + "\n" + string(actualContent))
		}
	}
}

const testToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiIwOGM0N2E4OC0yYzc5LTQyMGYtODEwNC02NWJkOWViYmU0MWUiLCJleHAiOjE1NDY1MDcyMzMsIm5iZiI6MCwiaWF0IjoxNTQ2NTA3MTczLCJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjgwMDEvYXV0aC9yZWFsbXMvbWFzdGVyIiwiYXVkIjoiZnJvbnRlbmQiLCJzdWIiOiJ0ZXN0T3duZXIiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJmcm9udGVuZCIsIm5vbmNlIjoiOTJjNDNjOTUtNzViMC00NmNmLTgwYWUtNDVkZDk3M2I0YjdmIiwiYXV0aF90aW1lIjoxNTQ2NTA3MDA5LCJzZXNzaW9uX3N0YXRlIjoiNWRmOTI4ZjQtMDhmMC00ZWI5LTliNjAtM2EwYWUyMmVmYzczIiwiYWNyIjoiMCIsImFsbG93ZWQtb3JpZ2lucyI6WyIqIl0sInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJ1c2VyIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsibWFzdGVyLXJlYWxtIjp7InJvbGVzIjpbInZpZXctcmVhbG0iLCJ2aWV3LWlkZW50aXR5LXByb3ZpZGVycyIsIm1hbmFnZS1pZGVudGl0eS1wcm92aWRlcnMiLCJpbXBlcnNvbmF0aW9uIiwiY3JlYXRlLWNsaWVudCIsIm1hbmFnZS11c2VycyIsInF1ZXJ5LXJlYWxtcyIsInZpZXctYXV0aG9yaXphdGlvbiIsInF1ZXJ5LWNsaWVudHMiLCJxdWVyeS11c2VycyIsIm1hbmFnZS1ldmVudHMiLCJtYW5hZ2UtcmVhbG0iLCJ2aWV3LWV2ZW50cyIsInZpZXctdXNlcnMiLCJ2aWV3LWNsaWVudHMiLCJtYW5hZ2UtYXV0aG9yaXphdGlvbiIsIm1hbmFnZS1jbGllbnRzIiwicXVlcnktZ3JvdXBzIl19LCJhY2NvdW50Ijp7InJvbGVzIjpbIm1hbmFnZS1hY2NvdW50IiwibWFuYWdlLWFjY291bnQtbGlua3MiLCJ2aWV3LXByb2ZpbGUiXX19LCJyb2xlcyI6WyJ1c2VyIl19.ykpuOmlpzj75ecSI6cHbCATIeY4qpyut2hMc1a67Ycg`

func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port, nil
}

func GetFreePort() (string, error) {
	temp, err := getFreePort()
	return strconv.Itoa(temp), err
}
