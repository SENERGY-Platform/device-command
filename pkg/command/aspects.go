/*
 * Copyright 2026 InfAI (CC SES)
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

package command

import (
	marshallermodel "github.com/SENERGY-Platform/marshaller/lib/marshaller/model"
)

//The single aspect of a command is deprecated in favor of an aspect list, following
//ContentVariable.AspectIds of the device-repository. It stays an alias for a list with one
//element: it is folded into the list where a command enters this package, so that a request
//that predates the lists is still handled and everything behind that boundary reads the
//list only.

// GetAspectIds returns the aspects the command asks for. A command names more than one
// aspect if the caller asked for more; a service then has to serve all of them, the way the
// device-repository reads a filter-criteria with several aspects.
func (this CommandMessage) GetAspectIds() []string {
	return marshallermodel.AspectIdsAlias(this.AspectId, this.AspectIds)
}

// SetAspectIds folds the deprecated AspectId into AspectIds and clears it, so that
// everything behind Command() and Batch() evaluates AspectIds only. Clearing matters for
// Batch: two elements naming the same aspect, one in the deprecated field and one in the
// list, are the same command and have to reach the same Hash().
func (this *CommandMessage) SetAspectIds() {
	this.AspectIds = this.GetAspectIds()
	this.AspectId = ""
}

// SetAspectIds folds the deprecated aspect of every element of the batch.
func (this BatchRequest) SetAspectIds() {
	for i := range this {
		this[i].SetAspectIds()
	}
}
