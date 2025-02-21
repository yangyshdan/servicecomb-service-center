/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util_test

import (
	"context"
	"testing"

	. "github.com/apache/servicecomb-service-center/datasource/etcd/util"
	"github.com/apache/servicecomb-service-center/pkg/util"
	pb "github.com/go-chassis/cari/discovery"
)

func TestFormatRevision(t *testing.T) {
	// null
	if x := FormatRevision(nil, nil); "da39a3ee5e6b4b0d3255bfef95601890afd80709" != x {
		t.Fatalf("TestFormatRevision failed, %s", x)
	}
	// 1.1,11.1,
	if x := FormatRevision([]int64{1, 11}, []int64{1, 1}); "87aa7d310290ff4f93248c0aed6870b928edf45a" != x {
		t.Fatalf("TestFormatRevision failed, %s", x)
	}
	// 1.11,1.1,
	if x := FormatRevision([]int64{1, 1}, []int64{11, 1}); "24675d196e3dea5be0c774cab281366640fc99ef" != x {
		t.Fatalf("TestFormatRevision failed, %s", x)
	}
}

func TestGetLeaseId(t *testing.T) {
	_, err := GetLeaseID(context.Background(), "", "", "")
	if err != nil {
		t.Fatalf(`GetLeaseID failed`)
	}
}

func TestGetInstance(t *testing.T) {
	_, err := GetInstance(context.Background(), "", "", "")
	if err != nil {
		t.Fatalf(`GetInstance failed`)
	}

	_, err = GetInstancesWithoutProperties(context.Background(), "", "")
	if err != nil {
		t.Fatalf(`GetInstancesWithoutProperties failed`)
	}

	_, err = QueryServiceInstancesKvs(context.Background(), "", 0)
	if err != nil {
		t.Fatalf(`QueryServiceInstancesKvs failed`)
	}
}

func TestDeleteServiceAllInstances(t *testing.T) {
	err := DeleteServiceAllInstances(context.Background(), "")
	if err != nil {
		t.Fatalf(`DeleteServiceAllInstances failed`)
	}
}

func TestParseEndpointValue(t *testing.T) {
	epv := ParseEndpointIndexValue([]byte("x/y"))
	if epv.ServiceID != "x" || epv.InstanceID != "y" {
		t.Fatalf(`ParseEndpointIndexValue failed`)
	}
}

func TestGetInstanceCountOfOneService(t *testing.T) {
	_, err := GetInstanceCountOfOneService(context.Background(), "", "")
	if err != nil {
		t.Fatalf(`GetInstanceCountOfOneService failed`)
	}
}

func TestUpdateInstance(t *testing.T) {
	err := UpdateInstance(util.WithNoCache(context.Background()), "", &pb.MicroServiceInstance{})
	if err == nil {
		t.Fatalf(`UpdateInstance CTX_NOCACHE failed`)
	}
}

func TestAppendFindResponse(t *testing.T) {
	ctx := context.Background()
	var (
		find              pb.FindInstancesResponse
		updatedResult     []*pb.FindResult
		notModifiedResult []int64
		failedResult      *pb.FindFailedResult
	)
	AppendFindResponse(ctx, 1, find.Response, find.Instances, &updatedResult, &notModifiedResult, &failedResult)
	if updatedResult == nil || notModifiedResult != nil || failedResult != nil {
		t.Fatal("TestAppendFindResponse failed")
	}
	if updatedResult[0].Index != 1 {
		t.Fatal("TestAppendFindResponse failed")
	}

	updatedResult = nil
	cloneCtx := util.WithResponseRev(ctx, "1")
	AppendFindResponse(cloneCtx, 1, find.Response, find.Instances, &updatedResult, &notModifiedResult, &failedResult)
	if updatedResult == nil || notModifiedResult != nil || failedResult != nil {
		t.Fatal("TestAppendFindResponse failed")
	}
	if updatedResult[0].Index != 1 || updatedResult[0].Rev != "1" {
		t.Fatal("TestAppendFindResponse failed")
	}

	updatedResult = nil
	cloneCtx = util.WithRequestRev(ctx, "1")
	cloneCtx = util.WithResponseRev(cloneCtx, "1")
	AppendFindResponse(cloneCtx, 1, find.Response, find.Instances, &updatedResult, &notModifiedResult, &failedResult)
	if updatedResult != nil || notModifiedResult == nil || failedResult != nil {
		t.Fatal("TestAppendFindResponse failed")
	}
	if notModifiedResult[0] != 1 {
		t.Fatal("TestAppendFindResponse failed")
	}

	notModifiedResult = nil
	find.Response = pb.CreateResponse(pb.ErrInternal, "test")
	AppendFindResponse(ctx, 1, find.Response, find.Instances, &updatedResult, &notModifiedResult, &failedResult)
	if updatedResult != nil || notModifiedResult != nil || failedResult == nil {
		t.Fatal("TestAppendFindResponse failed")
	}
	if failedResult.Error.Code != pb.ErrInternal {
		t.Fatal("TestAppendFindResponse failed")
	}
	find.Response = pb.CreateResponse(pb.ErrInvalidParams, "test")
	AppendFindResponse(ctx, 2, find.Response, find.Instances, &updatedResult, &notModifiedResult, &failedResult)
	if updatedResult != nil || notModifiedResult != nil || failedResult == nil {
		t.Fatal("TestAppendFindResponse failed")
	}
	if failedResult.Error.Code != pb.ErrInternal {
		t.Fatal("TestAppendFindResponse failed")
	}

	failedResult = nil
	find.Response = nil
	AppendFindResponse(ctx, 1, find.Response, find.Instances, &updatedResult, &notModifiedResult, &failedResult)
	AppendFindResponse(ctx, 2, find.Response, find.Instances, &updatedResult, &notModifiedResult, &failedResult)
	cloneCtx = util.WithRequestRev(ctx, "1")
	cloneCtx = util.WithResponseRev(cloneCtx, "1")
	AppendFindResponse(cloneCtx, 3, find.Response, find.Instances, &updatedResult, &notModifiedResult, &failedResult)
	AppendFindResponse(cloneCtx, 4, find.Response, find.Instances, &updatedResult, &notModifiedResult, &failedResult)
	find.Response = pb.CreateResponse(pb.ErrInternal, "test")
	AppendFindResponse(ctx, 5, find.Response, find.Instances, &updatedResult, &notModifiedResult, &failedResult)
	AppendFindResponse(ctx, 6, find.Response, find.Instances, &updatedResult, &notModifiedResult, &failedResult)
	if updatedResult == nil || notModifiedResult == nil || failedResult == nil {
		t.Fatal("TestAppendFindResponse failed")
	}
	if len(updatedResult) != 2 || len(notModifiedResult) != 2 || len(failedResult.Indexes) != 2 {
		t.Fatal("TestAppendFindResponse failed")
	}
}
