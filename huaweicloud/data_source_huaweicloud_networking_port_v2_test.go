package huaweicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2PortDataSource_basic(t *testing.T) {
	resourceName := "data.huaweicloud_networking_port.gw_port"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "all_fixed_ips.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "mac_address"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

func testAccNetworkingV2PortDataSource_basic() string {
	return fmt.Sprintf(`
data "huaweicloud_vpc_subnet" "mynet" {
  name = "subnet-default"
}

data "huaweicloud_networking_port" "gw_port" {
  network_id = data.huaweicloud_vpc_subnet.mynet.id
  fixed_ip   = data.huaweicloud_vpc_subnet.mynet.gateway_ip
}
`)
}
