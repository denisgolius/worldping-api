package api

import (
	"github.com/grafana/grafana/pkg/bus"
	"github.com/grafana/grafana/pkg/middleware"
	m "github.com/grafana/grafana/pkg/models"
	"strings"
)

func GetMonitorById(c *middleware.Context) {
	id := c.ParamsInt64(":id")

	query := m.GetMonitorByIdQuery{Id: id, OrgId: c.OrgId}
	query.IsRaintankAdmin = c.OrgRole == m.ROLE_RAINTANK_ADMIN

	err := bus.Dispatch(&query)
	if err != nil {
		c.JsonApiErr(404, "Monitor not found", nil)
		return
	}

	c.JSON(200, query.Result)
}

func GetMonitors(c *middleware.Context, query m.GetMonitorsQuery) {
	query.OrgId = c.OrgId
	query.IsRaintankAdmin = c.OrgRole == m.ROLE_RAINTANK_ADMIN

	if err := bus.Dispatch(&query); err != nil {
		c.JsonApiErr(500, "Failed to query monitors", err)
		return
	}
	c.JSON(200, query.Result)
}

func GetMonitorTypes(c *middleware.Context) {
	query := m.GetMonitorTypesQuery{}
	err := bus.Dispatch(&query)

	if err != nil {
		c.JsonApiErr(500, "Failed to query monitor_types", err)
		return
	}
	c.JSON(200, query.Result)
}

func DeleteMonitor(c *middleware.Context) {
	id := c.ParamsInt64(":id")

	cmd := &m.DeleteMonitorCommand{Id: id, OrgId: c.OrgId}

	err := bus.Dispatch(cmd)
	if err != nil {
		c.JsonApiErr(500, "Failed to delete monitor", err)
		return
	}

	c.JsonOK("monitor deleted")
}

func AddMonitor(c *middleware.Context) {
	cmd := m.AddMonitorCommand{}

	if !c.JsonBody(&cmd) {
		c.JsonApiErr(400, "Validation failed", nil)
		return
	}

	cmd.OrgId = c.OrgId
	if !c.IsGrafanaAdmin && strings.HasPrefix(strings.ToLower(cmd.Namespace), "public") {
		c.JsonApiErr(400, "Validation failed. namespace public is reserved.", nil)
		return
	}
	if cmd.Namespace == "" {
		cmd.Namespace = "network."
	}

	if err := bus.Dispatch(&cmd); err != nil {
		c.JsonApiErr(500, "Failed to add monitor", err)
		return
	}

	c.JSON(200, cmd.Result)
}

func UpdateMonitor(c *middleware.Context) {
	cmd := m.UpdateMonitorCommand{}

	if !c.JsonBody(&cmd) {
		c.JsonApiErr(400, "Validation failed", nil)
		return
	}

	cmd.OrgId = c.OrgId
	if !c.IsGrafanaAdmin && strings.HasPrefix(strings.ToLower(cmd.Namespace), "public") {
		c.JsonApiErr(400, "Validation failed. namespace public is reserved.", nil)
		return
	}

	err := bus.Dispatch(&cmd)
	if err != nil {
		c.JsonApiErr(500, "Failed to update monitor", err)
		return
	}

	c.JsonOK("Monitor updated")
}
