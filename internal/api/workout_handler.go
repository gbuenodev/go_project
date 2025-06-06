package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gbuenodev/goProject/internal/middleware"
	"github.com/gbuenodev/goProject/internal/store"
	"github.com/gbuenodev/goProject/internal/utils"
)

type WorkoutHandler struct {
	workoutStore store.WorkoutStore
	logger       *slog.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *slog.Logger) *WorkoutHandler {
	return &WorkoutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkoutHandler) HandleGetWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Error("ReadIDParam", "err", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout id"})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		if err == sql.ErrNoRows {
			wh.logger.Error("GetWorkoutByID", "err", err)
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
			return
		} else {
			wh.logger.Error("GetWorkoutByID", "err", err)
			utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to retrieve the workout"})
			return
		}
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": workout})
}

func (wh *WorkoutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Error("DecodingCreateWorkout", "err", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "must be logged in"})
		return
	}

	workout.UserID = currentUser.ID

	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Error("CreateWorkout", "err", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "failed to create workout"})
		return
	}

	utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"workout": createdWorkout})
}

func (wh *WorkoutHandler) HandleUpdateWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Error("ReadIDParam", "err", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout update id"})
		return
	}

	existingWorkout, err := wh.workoutStore.GetWorkoutByID(workoutID)
	if err != nil {
		wh.logger.Error("GetWorkoutByID", "err", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	} else if existingWorkout == nil {
		wh.logger.Error("GetWorkoutByID", "err", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	}

	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}

	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)
	if err != nil {
		wh.logger.Error("DecodingUpdateWorkout", "err", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid request"})
		return
	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		wh.logger.Debug("HandleUpdateWorkoutByID: no login detected")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "must be logged in"})
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			wh.logger.Debug("GetWorkoutOwner: the workout doesn't exist for the specified user")
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout does not exist"})
			return
		}

		wh.logger.Error("GetWorkoutOwner", "err", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if workoutOwner != currentUser.ID {
		wh.logger.Debug("UpdateWorkoutByID: user not allowed to update the specified workout")
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "not authorized to perform this action"})
		return
	}

	err = wh.workoutStore.UpdateWorkoutByID(existingWorkout)
	if err != nil {
		wh.logger.Error("UpdateWorkoutByID", "err", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, utils.Envelope{"workout": existingWorkout})
}

func (wh *WorkoutHandler) HandleDeleteWorkoutByID(w http.ResponseWriter, r *http.Request) {
	workoutID, err := utils.ReadIDParam(r)
	if err != nil {
		wh.logger.Error("ReadIDParam", "err", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "invalid workout delete id"})
		return
	}

	currentUser := middleware.GetUser(r)
	if currentUser == nil || currentUser == store.AnonymousUser {
		wh.logger.Debug("HandleDeleteWorkoutByID: no login detected")
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelope{"error": "must be logged in"})
		return
	}

	workoutOwner, err := wh.workoutStore.GetWorkoutOwner(workoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			wh.logger.Debug("GetWorkoutOwner: the workout doesn't exist for the specified user")
			utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout does not exist"})
			return
		}

		wh.logger.Error("GetWorkoutOwner", "err", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	if workoutOwner != currentUser.ID {
		wh.logger.Debug("DeleteWorkoutByID: user not allowed to update the specified workout")
		utils.WriteJSON(w, http.StatusForbidden, utils.Envelope{"error": "not authorized to perform this action"})
		return
	}

	err = wh.workoutStore.DeleteWorkoutByID(workoutID)
	if err == sql.ErrNoRows {
		wh.logger.Error("DeleteWorkoutByID", "err", err)
		utils.WriteJSON(w, http.StatusNotFound, utils.Envelope{"error": "workout not found"})
		return
	} else if err != nil {
		wh.logger.Error("DeleteWorkoutByID", "err", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelope{"error": "internal server error"})
		return
	}

	utils.WriteJSON(w, http.StatusNoContent, utils.Envelope{"message": "workout deleted"})
}
