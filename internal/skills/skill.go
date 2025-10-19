package skills

import (
	"context"
	"fmt"
	"sync"
)

// Skill represents a capability that can be invoked by the voice agent
type Skill interface {
	// Name returns the skill name
	Name() string

	// Description returns what the skill does
	Description() string

	// Execute runs the skill with given parameters
	Execute(ctx context.Context, params map[string]interface{}) (interface{}, error)

	// Schema returns the JSON schema for parameters
	Schema() map[string]interface{}
}

// Registry manages available skills
type Registry struct {
	skills map[string]Skill
	mu     sync.RWMutex
}

// NewRegistry creates a new skill registry
func NewRegistry() *Registry {
	return &Registry{
		skills: make(map[string]Skill),
	}
}

// Register registers a new skill
func (r *Registry) Register(skill Skill) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := skill.Name()
	if _, exists := r.skills[name]; exists {
		return fmt.Errorf("skill %s already registered", name)
	}

	r.skills[name] = skill
	return nil
}

// Get retrieves a skill by name
func (r *Registry) Get(name string) (Skill, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	skill, ok := r.skills[name]
	return skill, ok
}

// List returns all registered skills
func (r *Registry) List() []Skill {
	r.mu.RLock()
	defer r.mu.RUnlock()

	skills := make([]Skill, 0, len(r.skills))
	for _, skill := range r.skills {
		skills = append(skills, skill)
	}
	return skills
}

// Execute runs a skill by name
func (r *Registry) Execute(ctx context.Context, name string, params map[string]interface{}) (interface{}, error) {
	skill, ok := r.Get(name)
	if !ok {
		return nil, fmt.Errorf("skill %s not found", name)
	}

	return skill.Execute(ctx, params)
}

// --- Example Skills ---

// EchoSkill is a simple example skill
type EchoSkill struct{}

func (s *EchoSkill) Name() string {
	return "echo"
}

func (s *EchoSkill) Description() string {
	return "Echoes back the input message"
}

func (s *EchoSkill) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	message, ok := params["message"].(string)
	if !ok {
		return nil, fmt.Errorf("message parameter is required")
	}

	return map[string]interface{}{
		"echo": message,
	}, nil
}

func (s *EchoSkill) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"message": map[string]interface{}{
				"type":        "string",
				"description": "The message to echo back",
			},
		},
		"required": []string{"message"},
	}
}

// TimeSkill returns the current time
type TimeSkill struct{}

func (s *TimeSkill) Name() string {
	return "get_time"
}

func (s *TimeSkill) Description() string {
	return "Returns the current date and time"
}

func (s *TimeSkill) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	// In a real implementation, you might handle timezone parameter
	return map[string]interface{}{
		"current_time": fmt.Sprintf("%v", context.Background()),
	}, nil
}

func (s *TimeSkill) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}
}

// WeatherSkill is an example skill (stub)
type WeatherSkill struct{}

func (s *WeatherSkill) Name() string {
	return "get_weather"
}

func (s *WeatherSkill) Description() string {
	return "Gets weather information for a location"
}

func (s *WeatherSkill) Execute(ctx context.Context, params map[string]interface{}) (interface{}, error) {
	location, ok := params["location"].(string)
	if !ok {
		return nil, fmt.Errorf("location parameter is required")
	}

	// TODO: Integrate with real weather API
	return map[string]interface{}{
		"location":    location,
		"temperature": "72°F",
		"conditions":  "Sunny",
		"note":        "This is stub data - integrate with a real weather API",
	}, nil
}

func (s *WeatherSkill) Schema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "The city or location to get weather for",
			},
		},
		"required": []string{"location"},
	}
}

// InitDefaultSkills registers default skills
func InitDefaultSkills(registry *Registry) error {
	skills := []Skill{
		&EchoSkill{},
		&TimeSkill{},
		&WeatherSkill{},
	}

	for _, skill := range skills {
		if err := registry.Register(skill); err != nil {
			return err
		}
	}

	return nil
}
