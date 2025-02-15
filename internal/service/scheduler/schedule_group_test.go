package scheduler_test

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfscheduler "github.com/hashicorp/terraform-provider-aws/internal/service/scheduler"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccSchedulerScheduleGroup_basic(t *testing.T) {
	var scheduleGroup scheduler.GetScheduleGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_scheduler_schedule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.SchedulerEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SchedulerEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckScheduleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScheduleGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScheduleGroupExists(resourceName, &scheduleGroup),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "scheduler", regexp.MustCompile(regexp.QuoteMeta(`schedule-group/`+rName))),
					resource.TestCheckResourceAttrWith(resourceName, "creation_date", func(actual string) error {
						expect := scheduleGroup.CreationDate.Format(time.RFC3339)
						if actual != expect {
							return fmt.Errorf("expected value to be a formatted date")
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith(resourceName, "last_modification_date", func(actual string) error {
						expect := scheduleGroup.LastModificationDate.Format(time.RFC3339)
						if actual != expect {
							return fmt.Errorf("expected value to be a formatted date")
						}
						return nil
					}),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "state", "ACTIVE"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSchedulerScheduleGroup_disappears(t *testing.T) {
	var scheduleGroup scheduler.GetScheduleGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_scheduler_schedule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.SchedulerEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SchedulerEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckScheduleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScheduleGroupConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScheduleGroupExists(resourceName, &scheduleGroup),
					acctest.CheckResourceDisappears(acctest.Provider, tfscheduler.ResourceScheduleGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccSchedulerScheduleGroup_nameGenerated(t *testing.T) {
	var scheduleGroup scheduler.GetScheduleGroupOutput
	resourceName := "aws_scheduler_schedule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.SchedulerEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SchedulerEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckScheduleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScheduleGroupConfig_nameGenerated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScheduleGroupExists(resourceName, &scheduleGroup),
					acctest.CheckResourceAttrNameGenerated(resourceName, "name"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", resource.UniqueIdPrefix),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSchedulerScheduleGroup_namePrefix(t *testing.T) {
	var scheduleGroup scheduler.GetScheduleGroupOutput
	resourceName := "aws_scheduler_schedule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.SchedulerEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SchedulerEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckScheduleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScheduleGroupConfig_namePrefix("tf-acc-test-prefix-"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScheduleGroupExists(resourceName, &scheduleGroup),
					acctest.CheckResourceAttrNameFromPrefix(resourceName, "name", "tf-acc-test-prefix-"),
					resource.TestCheckResourceAttr(resourceName, "name_prefix", "tf-acc-test-prefix-"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccSchedulerScheduleGroup_tags(t *testing.T) {
	var scheduleGroup scheduler.GetScheduleGroupOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_scheduler_schedule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPartitionHasService(names.SchedulerEndpointID, t)
			testAccPreCheck(t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.SchedulerEndpointID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckScheduleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccScheduleGroupConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScheduleGroupExists(resourceName, &scheduleGroup),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccScheduleGroupConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScheduleGroupExists(resourceName, &scheduleGroup),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccScheduleGroupConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckScheduleGroupExists(resourceName, &scheduleGroup),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckScheduleGroupDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SchedulerClient
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_scheduler_schedule_group" {
			continue
		}

		_, err := conn.GetScheduleGroup(ctx, &scheduler.GetScheduleGroupInput{
			Name: aws.String(rs.Primary.ID),
		})
		if err != nil {
			var nfe *types.ResourceNotFoundException
			if errors.As(err, &nfe) {
				return nil
			}
			return err
		}

		return create.Error(names.Scheduler, create.ErrActionCheckingDestroyed, tfscheduler.ResNameScheduleGroup, rs.Primary.ID, errors.New("not destroyed"))
	}

	return nil
}

func testAccCheckScheduleGroupExists(name string, schedulegroup *scheduler.GetScheduleGroupOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.Scheduler, create.ErrActionCheckingExistence, tfscheduler.ResNameScheduleGroup, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.Scheduler, create.ErrActionCheckingExistence, tfscheduler.ResNameScheduleGroup, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).SchedulerClient
		ctx := context.Background()
		resp, err := conn.GetScheduleGroup(ctx, &scheduler.GetScheduleGroupInput{
			Name: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return create.Error(names.Scheduler, create.ErrActionCheckingExistence, tfscheduler.ResNameScheduleGroup, rs.Primary.ID, err)
		}

		*schedulegroup = *resp

		return nil
	}
}

func testAccPreCheck(t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).SchedulerClient
	ctx := context.Background()

	input := &scheduler.ListScheduleGroupsInput{}
	_, err := conn.ListScheduleGroups(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccScheduleGroupConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_scheduler_schedule_group" "test" {
  name = %[1]q
}
`, rName)
}

const testAccScheduleGroupConfig_nameGenerated = `
resource "aws_scheduler_schedule_group" "test" {}
`

func testAccScheduleGroupConfig_namePrefix(namePrefix string) string {
	return fmt.Sprintf(`
resource "aws_scheduler_schedule_group" "test" {
  name_prefix = %[1]q
}
`, namePrefix)
}

func testAccScheduleGroupConfig_tags1(rName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_scheduler_schedule_group" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1)
}

func testAccScheduleGroupConfig_tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_scheduler_schedule_group" "test" {
  name = %[1]q

  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2)
}
