# PRD: Multi-Profile Manager for Personalized Prompting

**Product Name**: Persona Profiles  
**Author**: Rahul Guha  
**Date**: Sep 2025

---

## 1. Objective

Allow users to create, edit, and switch between **multiple profiles** (personas). Each profile generates a personalized system prompt based on a free-text description (preferred entry) plus structured attributes.

---

## 2. Goals & Non-Goals

### Goals

- Users can create multiple profiles and switch between them.
- Encourage users to start with a **description** rather than filling forms.
- Automatically parse structured attributes from the description.
- Allow editing of parsed attributes before final save.
- Store all profiles and track the **active profile**.

### Non-Goals

- Not building a public persona/social feature.
- Not enforcing real identity — fictional or stylized personas are allowed.
- Not focusing on analytics (only personalization).

---

## 3. User Stories

1. As a **user**, I can see my profiles in a dropdown and select one as “active.”
2. As a **user**, I can create a new profile by writing a description.
3. As a **user**, I see structured fields extracted from my description and can adjust them.
4. As a **user**, I can save the profile and reuse/edit it later.
5. As a **developer**, I need each profile represented as a structured JSON for prompt injection.

---

## 4. User Flow

### Step 1: Profile Selection

- UI: **Dropdown of saved profiles** (shows profile names).
- **+ New Profile** button available.
- Selected profile becomes **active profile** for the session.

### Step 2: Create / Edit Profile

- Screen shows:
  - Large **Description Text Area** (primary entry).
  - Placeholder: _“I am a 25-year-old student in Dallas studying Computer Science and love history.”_
- Encourage user to write description (highlight: _“Best results come from writing a description.”_).

### Step 3: Auto-Generated Fields

- After description is entered → system parses and generates key-value fields:
  - Gender, Age, Location, Occupation, Interests, Education, Tone Preference, etc.
- UI shows description (unchanged) + structured attributes below.
- User can update/add/remove values.

### Step 4: Final Save

- User confirms → profile is saved.
- Redirect to **Profile Selection** screen with new profile visible.
- Saved profile becomes **active** by default.

---

## 5. Data Model

```json
{
  "profiles": [
    {
      "id": "uuid123",
      "name": "Student",
      "description": "I am a 25-year-old student in Dallas studying Computer Science and love history.",
      "attributes": {
        "gender": "Male",
        "age": 25,
        "location": {
          "city": "Dallas",
          "state": "Texas",
          "country": "USA",
          "postal_code": "75201"
        },
        "education_level": "Undergraduate",
        "occupation": "Student",
        "interests": ["Computer Science", "History"],
        "expertise_level": "Intermediate",
        "tone_preference": "Casual",
        "preferred_languages": ["English"],
        "intent": "To learn and explore new topics"
      }
    }
  ],
  "active_profile_id": "uuid123"
}
```
