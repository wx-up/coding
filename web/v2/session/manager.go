package session

import web "github.com/wx-up/coding/web/v2"

// Manager 为了用户体验，它不是必不可少的，用户完全可以自己拼凑 Store 和 Propagator
// 虽然在用户眼里它是一个核心结构，因为用户操作的直接是它，但是从设计者的角度，它就是一个胶水
type Manager struct {
	Store
	Propagator
}

// GetSession 获取 session
func (m *Manager) GetSession(ctx *web.Context) (Session, error) {
	if ctx.UserValues == nil {
		ctx.UserValues = make(map[string]any, 1)
	}

	// 从 request 获取到 sessionID
	id, err := m.Extract(ctx.Req.Context(), ctx.Req)
	if err != nil {
		return nil, err
	}

	// 尝试从缓存中获取
	value, ok := ctx.UserValues[id]
	if ok {
		return value.(Session), nil
	}

	// 目前 Get 是直接从 内存cache 中获取的，但是其他实现有可能是从 mysql 或者 redis 获取，所以缓存一下提高性能
	// 也有一种实现：即使是存储在 redis 中，但是 Get 的时候仍旧按照目前的设计从 内存cache 中获取
	// 当调用 session.Get 的时候再从 redis 中获取需要的 key
	session, err := m.Get(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}

	ctx.UserValues[id] = session
	return session, nil

}

// RefreshSession 刷新 session
func (m *Manager) RefreshSession(ctx *web.Context) (Session, error) {
	session, err := m.GetSession(ctx)
	if err != nil {
		return nil, err
	}

	// 刷新
	err = m.Refresh(ctx.Req.Context(), session.ID())
	if err != nil {
		return nil, err
	}

	// 重新写入 response 中
	err = m.Inject(ctx.Req.Context(), session.ID(), ctx.Resp)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// InitSession 初始化 session
func (m *Manager) InitSession(ctx *web.Context, id string) (Session, error) {
	session, err := m.Generate(ctx.Req.Context(), id)
	if err != nil {
		return nil, err
	}

	// 将 sessionID 写入到 response 中
	if err = m.Inject(ctx.Req.Context(), id, ctx.Resp); err != nil {
		return nil, err
	}

	return session, nil
}

// RemoveSession 删除 session
func (m *Manager) RemoveSession(ctx *web.Context) error {
	session, err := m.GetSession(ctx)
	if err != nil {
		return err
	}
	err = m.Store.Remove(ctx.Req.Context(), session.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(ctx.Req.Context(), ctx.Resp)
}
