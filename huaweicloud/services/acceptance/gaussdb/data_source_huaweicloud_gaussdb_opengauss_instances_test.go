package gaussdb

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
)

func TestAccOpenGaussInstancesDataSource_basic(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_gaussdb_opengauss_instances.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstancesDataSource_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenGaussInstancesDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.sharding_num", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.coordinator_num", "2"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.volume.0.size", "40"),
				),
			},
		},
	})
}

func TestAccOpenGaussInstancesDataSource_haModeCentralized(t *testing.T) {
	rName := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(5))
	dataSourceName := "data.huaweicloud_gaussdb_opengauss_instances.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenGaussInstancesDataSource_haModeCentralized(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenGaussInstancesDataSourceID(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "instances.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.replica_num", "3"),
					resource.TestCheckResourceAttr(dataSourceName, "instances.0.volume.0.size", "40"),
				),
			},
		},
	})
}

func testAccCheckOpenGaussInstancesDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("can't find GaussDB opengauss instances data source: %s ", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("the GaussDB opengauss data source ID not set ")
		}

		return nil
	}
}

func testAccOpenGaussInstancesDataSource_basic(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_instances" "test" {
  name = huaweicloud_gaussdb_opengauss_instance.test.name
  depends_on = [
    huaweicloud_gaussdb_opengauss_instance.test,
  ]
}
`, testAccOpenGaussInstance_basic(rName, fmt.Sprintf("%s@123", acctest.RandString(5)), 2))
}

func testAccOpenGaussInstancesDataSource_haModeCentralized(rName string) string {
	return fmt.Sprintf(`
%[1]s

data "huaweicloud_gaussdb_opengauss_instances" "test" {
  name = huaweicloud_gaussdb_opengauss_instance.test.name
  depends_on = [
    huaweicloud_gaussdb_opengauss_instance.test,
  ]
}
`, testAccOpenGaussInstance_haModeCentralized(rName, fmt.Sprintf("%s@123", acctest.RandString(5))))
}
