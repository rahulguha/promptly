package storage

import (
    "encoding/json"
    "errors"
    "os"
    "path/filepath"

    "github.com/google/uuid"
    "github.com/rahulguha/promptly/models"
)

type FileStorage struct {
    filePath string
}

func NewFileStorage(filePath string) *FileStorage {
    return &FileStorage{filePath: filePath}
}

func (fs *FileStorage) load() ([]models.Prompt, error) {
    if _, err := os.Stat(fs.filePath); errors.Is(err, os.ErrNotExist) {
        return []models.Prompt{}, nil
    }
    data, err := os.ReadFile(fs.filePath)
    if err != nil {
        return nil, err
    }
    var prompts []models.Prompt
    if err := json.Unmarshal(data, &prompts); err != nil {
        return nil, err
    }
    return prompts, nil
}

func (fs *FileStorage) save(prompts []models.Prompt) error {
    data, err := json.MarshalIndent(prompts, "", "  ")
    if err != nil {
        return err
    }
    os.MkdirAll(filepath.Dir(fs.filePath), 0755)
    return os.WriteFile(fs.filePath, data, 0644)
}

func (fs *FileStorage) SavePrompt(p models.Prompt) error {
    prompts, err := fs.load()
    if err != nil {
        return err
    }
    p.ID = uuid.New()
    prompts = append(prompts, p)
    return fs.save(prompts)
}

func (fs *FileStorage) GetPrompts() ([]models.Prompt, error) {
    return fs.load()
}

func (fs *FileStorage) SearchPrompt(userRole, llmRole string) ([]models.Prompt, error) {
    prompts, err := fs.load()
    if err != nil {
        return nil, err
    }
    var results []models.Prompt
    for _, pr := range prompts {
        if (userRole == "" || pr.UserRole == userRole) &&
           (llmRole == "" || pr.LLMRole == llmRole) {
            results = append(results, pr)
        }
    }
    return results, nil
}

func (fs *FileStorage) GetRoles(userRole string) (models.RoleResponse, error) {
    prompts, err := fs.load()
    if err != nil {
        return models.RoleResponse{}, err
    }
    userRoles := map[string]bool{}
    llmRoles := map[string]bool{}
    for _, pr := range prompts {
        userRoles[pr.UserRole] = true
        if pr.UserRole == userRole {
            llmRoles[pr.LLMRole] = true
        }
    }

    resp := models.RoleResponse{}
    if userRole == "" {
        for ur := range userRoles {
            resp.UserRoles = append(resp.UserRoles, ur)
        }
    } else {
        for lr := range llmRoles {
            resp.LLMRoles = append(resp.LLMRoles, lr)
        }
    }
    return resp, nil
}
