// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package wafregional_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/YakDriver/regexache"
	awstypes "github.com/aws/aws-sdk-go-v2/service/wafregional/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfwafregional "github.com/hashicorp/terraform-provider-aws/internal/service/wafregional"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func TestAccWAFRegionalRuleGroup_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var rule awstypes.Rule
	var group awstypes.RuleGroup
	var idx int

	ruleName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	groupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_rule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, names.WAFRegionalEndpointID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.WAFRegionalServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccRuleGroupConfig_basic(ruleName, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleExists(ctx, "aws_wafregional_rule.test", &rule),
					testAccCheckRuleGroupExists(ctx, resourceName, &group),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, "waf-regional", regexache.MustCompile(`rulegroup/.+`)),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, groupName),
					resource.TestCheckResourceAttr(resourceName, "activated_rule.#", acctest.CtOne),
					resource.TestCheckResourceAttr(resourceName, names.AttrMetricName, groupName),
					computeActivatedRuleWithRuleId(&rule, "COUNT", 50, &idx),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "activated_rule.*", map[string]string{
						"action.0.type":    "COUNT",
						names.AttrPriority: "50",
						names.AttrType:     string(awstypes.WafRuleTypeRegular),
					}),
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

func TestAccWAFRegionalRuleGroup_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var rule awstypes.Rule
	var group awstypes.RuleGroup

	ruleName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	groupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_rule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, names.WAFRegionalEndpointID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.WAFRegionalServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccRuleGroupConfig_tags1(ruleName, groupName, acctest.CtKey1, acctest.CtValue1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleExists(ctx, "aws_wafregional_rule.test", &rule),
					testAccCheckRuleGroupExists(ctx, resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.CtOne),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey1, acctest.CtValue1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccRuleGroupConfig_tags2(ruleName, groupName, acctest.CtKey1, acctest.CtValue1Updated, acctest.CtKey2, acctest.CtValue2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleExists(ctx, "aws_wafregional_rule.test", &rule),
					testAccCheckRuleGroupExists(ctx, resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.CtTwo),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey1, acctest.CtValue1Updated),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey2, acctest.CtValue2),
				),
			},
			{
				Config: testAccRuleGroupConfig_tags1(ruleName, groupName, acctest.CtKey2, acctest.CtValue2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleExists(ctx, "aws_wafregional_rule.test", &rule),
					testAccCheckRuleGroupExists(ctx, resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsPercent, acctest.CtOne),
					resource.TestCheckResourceAttr(resourceName, acctest.CtTagsKey2, acctest.CtValue2),
				),
			},
		},
	})
}

func TestAccWAFRegionalRuleGroup_changeNameForceNew(t *testing.T) {
	ctx := acctest.Context(t)
	var before, after awstypes.RuleGroup

	ruleName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	groupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	newGroupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_rule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, names.WAFRegionalEndpointID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.WAFRegionalServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccRuleGroupConfig_basic(ruleName, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleGroupExists(ctx, resourceName, &before),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, groupName),
					resource.TestCheckResourceAttr(resourceName, "activated_rule.#", acctest.CtOne),
					resource.TestCheckResourceAttr(resourceName, names.AttrMetricName, groupName),
				),
			},
			{
				Config: testAccRuleGroupConfig_basic(ruleName, newGroupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleGroupExists(ctx, resourceName, &after),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, newGroupName),
					resource.TestCheckResourceAttr(resourceName, "activated_rule.#", acctest.CtOne),
					resource.TestCheckResourceAttr(resourceName, names.AttrMetricName, newGroupName),
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

func TestAccWAFRegionalRuleGroup_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var group awstypes.RuleGroup
	ruleName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	groupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_rule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, names.WAFRegionalEndpointID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.WAFRegionalServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccRuleGroupConfig_basic(ruleName, groupName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRuleGroupExists(ctx, resourceName, &group),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfwafregional.ResourceRuleGroup(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccWAFRegionalRuleGroup_changeActivatedRules(t *testing.T) {
	ctx := acctest.Context(t)
	var rule0, rule1, rule2, rule3 awstypes.Rule
	var groupBefore, groupAfter awstypes.RuleGroup
	var idx0, idx1, idx2, idx3 int

	groupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	ruleName1 := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	ruleName2 := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	ruleName3 := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_rule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, names.WAFRegionalEndpointID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.WAFRegionalServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccRuleGroupConfig_basic(ruleName1, groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRuleExists(ctx, "aws_wafregional_rule.test", &rule0),
					testAccCheckRuleGroupExists(ctx, resourceName, &groupBefore),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, groupName),
					resource.TestCheckResourceAttr(resourceName, "activated_rule.#", acctest.CtOne),
					computeActivatedRuleWithRuleId(&rule0, "COUNT", 50, &idx0),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "activated_rule.*", map[string]string{
						"action.0.type":    "COUNT",
						names.AttrPriority: "50",
						names.AttrType:     string(awstypes.WafRuleTypeRegular),
					}),
				),
			},
			{
				Config: testAccRuleGroupConfig_changeActivateds(ruleName1, ruleName2, ruleName3, groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, names.AttrName, groupName),
					resource.TestCheckResourceAttr(resourceName, "activated_rule.#", acctest.CtThree),
					testAccCheckRuleGroupExists(ctx, resourceName, &groupAfter),

					testAccCheckRuleExists(ctx, "aws_wafregional_rule.test", &rule1),
					computeActivatedRuleWithRuleId(&rule1, "BLOCK", 10, &idx1),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "activated_rule.*", map[string]string{
						"action.0.type":    "BLOCK",
						names.AttrPriority: "10",
						names.AttrType:     string(awstypes.WafRuleTypeRegular),
					}),

					testAccCheckRuleExists(ctx, "aws_wafregional_rule.test2", &rule2),
					computeActivatedRuleWithRuleId(&rule2, "COUNT", 1, &idx2),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "activated_rule.*", map[string]string{
						"action.0.type":    "COUNT",
						names.AttrPriority: acctest.CtOne,
						names.AttrType:     string(awstypes.WafRuleTypeRegular),
					}),

					testAccCheckRuleExists(ctx, "aws_wafregional_rule.test3", &rule3),
					computeActivatedRuleWithRuleId(&rule3, "BLOCK", 15, &idx3),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "activated_rule.*", map[string]string{
						"action.0.type":    "BLOCK",
						names.AttrPriority: "15",
						names.AttrType:     string(awstypes.WafRuleTypeRegular),
					}),
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

func TestAccWAFRegionalRuleGroup_noActivatedRules(t *testing.T) {
	ctx := acctest.Context(t)
	var group awstypes.RuleGroup
	groupName := fmt.Sprintf("tfacc%s", sdkacctest.RandString(5))
	resourceName := "aws_wafregional_rule_group.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckPartitionHasService(t, names.WAFRegionalEndpointID) },
		ErrorCheck:               acctest.ErrorCheck(t, names.WAFRegionalServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckRuleGroupDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccRuleGroupConfig_noActivateds(groupName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckRuleGroupExists(ctx, resourceName, &group),
					resource.TestCheckResourceAttr(resourceName, names.AttrName, groupName),
					resource.TestCheckResourceAttr(resourceName, "activated_rule.#", acctest.CtZero),
				),
			},
		},
	})
}

func testAccCheckRuleGroupDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_wafregional_rule_group" {
				continue
			}

			conn := acctest.Provider.Meta().(*conns.AWSClient).WAFRegionalClient(ctx)

			_, err := tfwafregional.FindRuleGroupByID(ctx, conn, rs.Primary.ID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("WAF Regional Rule Group %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckRuleGroupExists(ctx context.Context, n string, v *awstypes.RuleGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).WAFRegionalClient(ctx)

		output, err := tfwafregional.FindRuleGroupByID(ctx, conn, rs.Primary.ID)

		if err != nil {
			return err
		}

		*v = *output

		return nil
	}
}

func testAccRuleGroupConfig_basic(ruleName, groupName string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_rule" "test" {
  name        = %[1]q
  metric_name = %[1]q
}

resource "aws_wafregional_rule_group" "test" {
  name        = %[2]q
  metric_name = %[2]q

  activated_rule {
    action {
      type = "COUNT"
    }

    priority = 50
    rule_id  = aws_wafregional_rule.test.id
  }
}
`, ruleName, groupName)
}

func testAccRuleGroupConfig_tags1(ruleName, groupName, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_rule" "test" {
  name        = %[1]q
  metric_name = %[1]q
}

resource "aws_wafregional_rule_group" "test" {
  name        = %[2]q
  metric_name = %[2]q

  activated_rule {
    action {
      type = "COUNT"
    }

    priority = 50
    rule_id  = aws_wafregional_rule.test.id
  }

  tags = {
    %[3]q = %[4]q
  }
}
`, ruleName, groupName, tagKey1, tagValue1)
}

func testAccRuleGroupConfig_tags2(ruleName, groupName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_rule" "test" {
  name        = %[1]q
  metric_name = %[1]q
}

resource "aws_wafregional_rule_group" "test" {
  name        = %[2]q
  metric_name = %[2]q

  activated_rule {
    action {
      type = "COUNT"
    }

    priority = 50
    rule_id  = aws_wafregional_rule.test.id
  }

  tags = {
    %[3]q = %[4]q
    %[5]q = %[6]q
  }
}
`, ruleName, groupName, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccRuleGroupConfig_changeActivateds(ruleName1, ruleName2, ruleName3, groupName string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_rule" "test" {
  name        = %[1]q
  metric_name = %[1]q
}

resource "aws_wafregional_rule" "test2" {
  name        = %[2]q
  metric_name = %[2]q
}

resource "aws_wafregional_rule" "test3" {
  name        = %[3]q
  metric_name = %[3]q
}

resource "aws_wafregional_rule_group" "test" {
  name        = %[4]q
  metric_name = %[4]q

  activated_rule {
    action {
      type = "BLOCK"
    }

    priority = 10
    rule_id  = aws_wafregional_rule.test.id
  }

  activated_rule {
    action {
      type = "COUNT"
    }

    priority = 1
    rule_id  = aws_wafregional_rule.test2.id
  }

  activated_rule {
    action {
      type = "BLOCK"
    }

    priority = 15
    rule_id  = aws_wafregional_rule.test3.id
  }
}
`, ruleName1, ruleName2, ruleName3, groupName)
}

func testAccRuleGroupConfig_noActivateds(groupName string) string {
	return fmt.Sprintf(`
resource "aws_wafregional_rule_group" "test" {
  name        = %[1]q
  metric_name = %[1]q
}
`, groupName)
}

// computeActivatedRuleWithRuleId calculates index
// which isn't static because ruleId is generated as part of the test
func computeActivatedRuleWithRuleId(rule *awstypes.Rule, actionType string, priority int, idx *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ruleResource := tfwafregional.ResourceRuleGroup().SchemaMap()["activated_rule"].Elem.(*schema.Resource)

		m := map[string]interface{}{
			names.AttrAction: []interface{}{
				map[string]interface{}{
					names.AttrType: actionType,
				},
			},
			names.AttrPriority: priority,
			"rule_id":          *rule.RuleId,
			names.AttrType:     string(awstypes.WafRuleTypeRegular),
		}

		f := schema.HashResource(ruleResource)
		*idx = f(m)

		return nil
	}
}
